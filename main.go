package main

import(
	//"os"
	"log"
	"net"
	"fmt"
	"errors"
	"strings"
	"strconv"
	"ircd/m/user"
	"ircd/m/message"
	"ircd/m/server"
	"ircd/m/reply"
	"ircd/m/channel"
	"ircd/m/format"
	"ircd/m/input"
)

var(
	srv server.Server
	// Conection is closed error
)

func
HandleNick(a server.HndlArg) error {
	newNick := a.Msg.Args[1]

	_, nickExists := srv.Users[newNick]
	if nickExists {
		log.Printf("Nick '%s' is already taken\n", newNick)
		return message.Send(a.Usr, message.Message{srv.Host,
			[]string{reply.FormatNum(reply.ERR_NICKNAMEINUSE), "Nickname is already in use."},
		})
	}

	// Delete old user.
	if a.Usr.Nick != "" {
		delete(srv.Users, a.Usr.Nick)
	}

	log.Printf("Set nick of '%s' to '%s'\n", a.Usr.Nick, newNick)
	a.Usr.Nick = newNick
	srv.Users[newNick] = a.Usr
	return nil
}

func
HandlePass(arg server.HndlArg) error {
	return nil
}

func
HandleSquit(arg server.HndlArg) error {
	return nil
}

func
HandleUser(a server.HndlArg) error {
	user, mode, _, info :=
		a.Msg.Args[1], a.Msg.Args[2],
		a.Msg.Args[3], a.Msg.Args[4]
	
	a.Usr.User = user
	if v, err := strconv.Atoi(mode) ; err != nil {
		a.Usr.Mode = 0
	} else {
		a.Usr.Mode = v
	}
	a.Usr.Info = info
	
	return nil
}

func
HandleJoin(a server.HndlArg) error {
	chanStr := a.Msg.Args[1]
	chanNames := strings.Split(chanStr, channel.NamDel)
	for _, v := range chanNames {
		// Skip channel names without prefixes.
		if format.HasAnyOfPrefixes(v, channel.NamPre) == "" {
			continue	
		}

		ch, ok := srv.Chans[v]
		// Create new channel if does not exist.
		if !ok {
			srv.Chans[v] = &channel.Channel{ make(map[string]*user.User)}
			ch = srv.Chans[v]
		}

		ch.Users[a.Usr.Nick] = a.Usr
	}
	return nil
}

func
HandlePart(arg server.HndlArg) error {
	return nil
}

func
HandlePrivMsg(a server.HndlArg) error {
	var recvs []*user.User
	alltos := a.Msg.Args[1]
	msgstr := a.Msg.Args[2]
	names := strings.Split(alltos, ",")

	// Getting list of receivers.
	for _, to := range names {
		pref := format.HasAnyOfPrefixes(to, channel.NamPre)
		if pref != "" { // For channels.
			ch, ok := srv.Chans[to]
			if ok {
				for  _, v := range ch.Users {
					recvs = append(recvs, v)
				}
			}
		} else { // For exact user.
			usr, ok := srv.Users[to]
			if !ok {
				continue
			}
			recvs = append(recvs, usr)
		}
	}
	
	if len(recvs) == 0 {
		return message.Send(a.Usr,
			message.Message{
				a.Msg.Src,
				[]string{
					reply.FormatNum(reply.ERR_NORECIPIENT),
					fmt.Sprintf("No recipient given (%s)", a.Msg.Args[0]),
				},
			},
		)
	}

	// Sending to every of them.
	for _, u := range recvs {
		log.Printf("Sending private message to '%s'\n", a.Usr.Nick)
		err := message.Send(
			u,
			message.Message{u.FullSrc(), []string{a.Msg.Args[0], a.Usr.Nick, msgstr}})
		if err == input.CIC {
			CleanUpUser(u)
		}
	}
	return nil
}
func
HandleOper(arg server.HndlArg) error {
	return nil
}

func
HandleMode(arg server.HndlArg) error {
	return nil
}


func
HandleMessage(usr *user.User, msg message.Message) error {
	cmd, ok := srv.Cmds[msg.Args[0]]
	if !ok {
		return errors.New("No such command")
	}

	if cmd.Nargs > len(msg.Args) - 1 {
		log.Printf("Not enough arguments for '%s'\n", msg.Args[0])
		return nil
	}

	log.Printf("Handling '%s'", msg.Args[0])
	err := cmd.Hndl(server.HndlArg{usr, msg})
	if err != nil {
		return err
	}

	return nil
}

func
CleanUpUser(usr *user.User){
	/* We do not care if it is closed already. */
	usr.Conn.Close()

	if usr.Nick != "" {
		delete(srv.Users, usr.Nick)
	}
	log.Printf("Done cleaning '%s' user\n", usr.Nick)
}

func
HandleConn(conn net.Conn) {
	/* New user. Not in server lists yet though. */
	usr := user.User{Conn: conn}
	for {
		msg, err := message.Read(conn)
		/* If connection was closed in other thread 
			then also CIC will be returned*/
		if err == input.CIC {
			CleanUpUser(&usr)
			return;
		}
		fmt.Printf("'%s' %v\n", msg.Src, msg.Args)

		/* Handling includes writing replies,
			so CIC is checked when writing. */
		err = HandleMessage(&usr, msg)
		if err == input.CIC {
			CleanUpUser(&usr)
			return
		}
	}
}

func
main() {
	host := "localhost"
	port := 6667
	addrStr := host+":"+strconv.Itoa(port)
	fmt.Println(addrStr)
	ln, err := net.Listen("tcp", addrStr)
	if err != nil {
		log.Fatal(err)
	}

	srv = server.Server{Host: host, Port: port,
		Ln: ln,
		Users: make(map[string]*user.User),
		Chans: make(map[string]*channel.Channel),
		Cmds: server.Commands {
			"NICK":{ 1, HandleNick},
			"PASS":{ 1, HandlePass},
			"USER":{ 3, HandleUser},

			"JOIN":{ 1, HandleJoin},
			"PART":{ 1, HandlePart},
			"PRIVMSG":{ 2, HandlePrivMsg},
			"OPER":{ 0, HandleOper},
			"MODE":{ 0, HandleMode},
		},
	}

	for {
		conn, err := srv.Ln.Accept()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(conn.RemoteAddr())
		go HandleConn(conn);
	}
}
