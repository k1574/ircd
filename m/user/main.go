package user

import(
	"net"
)

/* Stick connection and nick together. */
type User struct {
	conn net.Conn
	nick, user, host, info string
	mode int
}

var(
	MaxNikLen = 9
)

func
(u *User)CleanUp(){
	usr.conn.close()
}

func
(u *User)FullSrc() string {
	return u.nick+"!"+u.user+"@"+u.host
}