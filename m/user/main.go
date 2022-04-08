package user

import(
	"net"
)

/* Stick connection and nick together. */
type User struct {
	Conn net.Conn
	Nick, User, Host, Info string
	Mode int
}

var(
	MaxNikLen = 9
)

func
(u *User)CleanUp(){
	u.Conn.Close()
}

func
(u *User)FullSrc() string {
	return u.Nick+"!"+u.User+"@"+u.Host
}