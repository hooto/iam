package main

import (
    "../deps/fsnotify"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "sync"
    "time"
)

var (
    cmd         *exec.Cmd
    lock        sync.Mutex
    paths       []string
    buildPeriod time.Time
    exts        []string
    apppath     string
    appname     string
    mainFiles   []string
    eventTime   = make(map[string]int64)
)

func main() {

    appname = "bin/lessids-server"
    apppath = "src"

    exts = []string{"go", "tpl", "json", "js", "css"}

    mainFiles = []string{"main.go"}
    for k, v := range mainFiles {
        mainFiles[k] = apppath + "/" + v
    }

    go Start(appname)

    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        fmt.Println("[ERRO] Fail to create new Watcher[ %s ]", err)
        os.Exit(2)
    }

    go func() {
        for {
            select {
            case e := <-watcher.Event:

                isbuild := false

                if buildPeriod.Add(1 * time.Second).After(time.Now()) {
                    continue
                }
                buildPeriod = time.Now()

                for _, ext := range exts {

                    if strings.HasSuffix(strings.ToLower(e.Name), "."+ext) {
                        fmt.Println("[INFO]", e)
                        isbuild = true
                    }
                }

                if isbuild {

                    mt := getFileModTime(e.Name)
                    if t := eventTime[e.Name]; mt == t {
                        fmt.Printf("[SKIP] %s\n", e.String())
                        isbuild = false
                    }

                    eventTime[e.Name] = mt
                }

                if isbuild {
                    Autobuild(mainFiles)
                }

            case err := <-watcher.Error:
                fmt.Println("[WARN] %s", err.Error())
            }
        }
    }()

    WalkFunc := func(path string, info os.FileInfo, err error) error {

        if err != nil || !info.IsDir() {
            return nil
        }

        watcher.Watch(path)
        return nil
    }
    _ = filepath.Walk(apppath, WalkFunc)

    for {
        time.Sleep(1e9)
    }
}

func getFileModTime(path string) int64 {

    f, err := os.Open(path)
    if err != nil {
        return time.Now().Unix()
    }
    defer f.Close()

    fi, err := f.Stat()
    if err != nil {
        return time.Now().Unix()
    }

    return fi.ModTime().Unix()
}

func Autobuild(files []string) {

    lock.Lock()
    defer lock.Unlock()

    var err error

    fmt.Println("[INFO] Start building...")

    args := []string{"build", "-o", appname}
    args = append(args, files...)

    bcmd := exec.Command("/usr/local/go/bin/go", args...)
    bcmd.Stdout = os.Stdout
    bcmd.Stderr = os.Stderr
    err = bcmd.Run()

    if err != nil {
        fmt.Println("[ERRO] Build failed")
        return
    }

    Restart(appname)
}

func Stop() {

    fmt.Println("[INFO] Stopping %s ...", appname)

    defer func() {
        if e := recover(); e != nil {
            fmt.Println("Kill.recover -> ", e)
        }
    }()

    if cmd != nil && cmd.Process != nil {
        err := cmd.Process.Kill()
        if err != nil {
            fmt.Println("Kill -> ", err)
        }
    }
}

func Restart(appname string) {
    Stop()
    go Start(appname)
}

func Start(appname string) {

    fmt.Printf("[INFO] Restarting %s ...\n", appname)
    if strings.Index(appname, "./") == -1 {
        appname = "./" + appname
    }

    cmd = exec.Command(appname)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Args = append([]string{appname})

    go cmd.Run()
    fmt.Printf("[INFO] %s is running...\n", appname)
}
