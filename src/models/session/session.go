package session

import (
    "../../../../deps/lessgo/pagelet"
    "sync"
    "time"
)

var (
    session Session = map[string]SessionItem{}
    locker  sync.Mutex
)

type SessionItem struct {
    UserID string
    Expire time.Time
}

// A signed cookie (and thus limited to 4kb in size).
// Restriction: Keys may not have a colon in them.
type Session map[string]SessionItem

func Set(key, userid, expire string) {

    locker.Lock()
    defer locker.Unlock()

    if _, ok := session[key]; ok {
        return
    }

    taf, _ := time.ParseDuration("+" + expire + "s")
    item := SessionItem{
        UserID: userid,
        Expire: time.Now().Add(taf),
    }

    session[key] = item
}

func IsLogin(r *pagelet.Request) bool {

    if cookie, err := r.Request.Cookie("access_token_lessfly"); err == nil {

        if _, ok := session[cookie.Value]; ok {
            return true
        }
    }

    return false
}
