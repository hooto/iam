package session

import (
	"strings"
	"sync"
	"time"

	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
)

var (
	locker sync.Mutex
)

type Session struct {
	Uid     string
	Uname   string
	Expired time.Time
}

func IsLogin(r *httpsrv.Request) bool {

	sess := GetSession(r)
	if sess.Uid != "" {
		return true
	}

	return false
}

func GetSession(r *httpsrv.Request) (sess Session) {

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

	sess.Uid = rsu[0].Field("uid").String()
	sess.Uname = rsu[0].Field("uname").String()

	return
}
