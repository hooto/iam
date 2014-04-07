package session

import (
    "../../../deps/lessgo/data/rdc"
    "../../../deps/lessgo/pagelet"
    "fmt"
    "strconv"
    "strings"
    "sync"
    "time"
)

var (
    locker sync.Mutex
)

type Session struct {
    Uid     uint32
    Uname   string
    Expired time.Time
}

func IsLogin(r *pagelet.Request) bool {

    sess := GetSession(r)
    if sess.Uid > 0 {
        return true
    }
    return false
}

func GetSession(r *pagelet.Request) (sess Session) {

    cookie, err := r.Request.Cookie("access_token")
    if err != nil {
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        return
    }

    q := rdc.NewQuerySet().From("ids_sessions").Limit(1)
    q.Where.And("token", cookie.Value)
    rsu, err := dcn.Query(q)
    if err == nil && len(rsu) == 0 {
        return
    }

    sess.Expired = rsu[0]["expired"].(time.Time)
    if sess.Expired.Before(time.Now()) {
        fmt.Println("Expired")
        return
    }

    addr := "0.0.0.0"
    if addridx := strings.Index(r.RemoteAddr, ":"); addridx > 0 {
        addr = r.RemoteAddr[:addridx]
    }
    if addr != rsu[0]["source"].(string) {
        fmt.Println("source")
        return
    }

    uid, _ := strconv.Atoi(fmt.Sprintf("%v", rsu[0]["uid"]))
    sess.Uid = uint32(uid)
    sess.Uname = rsu[0]["uname"].(string)

    return
}
