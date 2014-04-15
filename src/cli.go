package main

import (
    "../deps/lessgo/pass"
    "../deps/lessgo/utils"
    "../deps/lessgo/data/rdc"
    "./conf"
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "os/exec"
    "os/user"
    "regexp"
    "strings"
    "time"
)

var err error
var cfg conf.Config
var emailPattern = regexp.MustCompile("^[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,10})$")

var flagPrefix = flag.String("prefix", "", "the prefix folder path")
var flagUserSet = flag.String("userset", "", "the username")
var flagUserDel = flag.String("userdel", "", "the username")

func main() {

    //
    if u, err := user.Current(); err != nil || u.Uid != "0" {
        //log.Fatal("Permission Denied : must be run as root")
    }

    //
    flag.Parse()
    if cfg, err = conf.NewConfig(*flagPrefix); err != nil {
        log.Fatal(err)
    }

    if cn, err := cfg.DatabaseInstance(); err == nil {
        rdc.InstanceRegister("def", cn)
    } else {
        log.Fatal(err)
    }

    if *flagUserSet != "" {
        cmdUserSet()
    } else if *flagUserDel != "" {
        //cmdUserDel()
    } else {
        fmt.Println("No Command Found")
    }
}

func cmdUserSet() {

    defer func() {
        exec.Command("stty", "-F", "/dev/tty", "-cbreak").Run()
    }()

    // disable input buffering
    exec.Command("stty", "-F", "/dev/tty", "cbreak").Run()
    // do not display entered characters on the screen
    //exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

    //
    if *flagUserSet == "" {
        fmt.Println("Email can not be null")
        os.Exit(1)
    }
    email := strings.ToLower(strings.TrimSpace(*flagUserSet))
    if matched := emailPattern.MatchString(email); !matched {
        fmt.Println("Email is not valid")
        os.Exit(1)
    }

    //
    prompt := "\rEnter new password: "
    fmt.Printf("Setting password for %s\n%s", email, prompt)
    reader := bufio.NewReaderSize(os.Stdin, 1)
    passwd := ""
    for {

        c, _ := reader.ReadByte()
        if c == '\n' {
            break
        }

        passwd += string(c)

        prompt += "*"
        fmt.Print(prompt)
    }
    if len(passwd) < 12 || len(passwd) > 30 {
        fmt.Println("Password must be between 12 and 30 characters long")
        os.Exit(1)
    }
    hash, _ := pass.HashDefault(passwd)
    //fmt.Println(hash)

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        log.Fatal("Internal Server Error")
    }

    q := rdc.NewQuerySet().From("ids_login").Limit(1)
    q.Where.And("email", email)
    rsu, err := dcn.Query(q)
    if err == nil && len(rsu) == 1 {
        fmt.Println("The `Email` already exists, please choose another one")
        os.Exit(1)
    }

    uname := utils.StringNewRand36(8)
    item := map[string]interface{}{
        "uname":   uname,
        "email":   email,
        "pass":    hash,
        "name":    uname,
        "status":  1,
        "created": time.Now().Format(time.RFC3339), // TODO
        "updated": time.Now().Format(time.RFC3339), // TODO
    }
    if err := dcn.Insert("ids_login", item); err != nil {
        fmt.Println("Internal Server Error: Can not write to database")
        os.Exit(1)
    }

    //
    fmt.Println("Password updated successfully")
}

/*
func cmdUserDel() {

    //
    if *flagUserDel == "" {
        log.Fatal("Username can not be null")
    }
    email := strings.ToLower(*flagUserDel)



    fmt.Println("User deleted successfully")
}
*/
