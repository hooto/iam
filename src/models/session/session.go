package session

import (
    "../../../deps/lessgo/data/rdo"
    "../../../deps/lessgo/data/rdo/base"
    "../../../deps/lessgo/pagelet"
    "strings"
    "sync"
    "time"
)

var (
    locker sync.Mutex
)

type Session struct {
    Uid     uint32
    Uuid    string
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

    token, err := r.Request.Cookie("access_token")
    if err != nil {
        return
    }

    dcn, err := rdo.ClientPull("def")
    if err != nil {
        return
    }

    q := base.NewQuerySet().From("ids_sessions").Limit(1)
    q.Where.And("token", token.Value)
    rsu, err := dcn.Base.Query(q)
    if err == nil && len(rsu) == 0 {
        return
    }

    sess.Expired = rsu[0].Field("expired").TimeParse("datetime")
    if sess.Expired.Before(time.Now()) {
        return
    }

    addr := "0.0.0.0"
    if addridx := strings.Index(r.RemoteAddr, ":"); addridx > 0 {
        addr = r.RemoteAddr[:addridx]
    }
    if addr != rsu[0].Field("source").String() {
        return
    }

    sess.Uid = uint32(rsu[0].Field("uid").Int())
    sess.Uuid = rsu[0].Field("uuid").String()
    sess.Uname = rsu[0].Field("uname").String()

    return
}
