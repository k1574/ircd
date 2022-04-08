package server

import(
	"net"
	"ircd/m/message"
	"ircd/m/user"
	"ircd/m/channel"
)


/* Arguments for handlers. */
type HndlArg struct {
	Usr *user.User
	Msg message.Message
}

type Commands map[string]struct {
	Nargs int
	Hndl func(arg HndlArg) error
}

/* The main struct for all the project. */
type Server struct {
	Host string
	Port int
	Ln net.Listener
	Users map[string]*user.User
	Chans map[string]*channel.Channel
	Cmds Commands
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
}