package server

import(
	"net"
	"ircd/m/user"
	"ircd/m/channel"
)


/* Arguments for handlers. */
type HndlArg struct {
	usr *user.User
	msg message.Message
}

type Commands map[string]struct {
	nargs int
	hndl func(arg HndlArg)
}

/* The main struct for all the project. */
type Server struct {
	host string
	port int
	ln net.Listener
	users map[string]*user.User
	chans map[string]*channel.Channel
	cmds Commands
}

func
(srv *Server)UserExists(nick string) bool {
	return false
}

func
(srv *Server)RemoveUser(nick string) {
}

func
(srv *Server)AddChan(name string) {
}

func
(srv *Server)RemoveChan(name string) {
	srv
}