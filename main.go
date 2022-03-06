package main

import(
	//"os"
	"log"
	"net"
	"fmt"
	"io"
	"errors"
	"bufio"
	"strings"
	"strconv"
	"ircd/m/user"
	"ircd/m/message"
	"ircd/m/server"
)



var(
	srv server.Server
	// Conection is closed error
	CIC = errors.New("connection is closed")

	srv.cmds = server.Commands {
		"NICK":{ 1, HandleNick},
		"PASS":{ 1, HandlePass},
		"USER":{ 3, HandleUser},

		"JOIN":{ 1, HandleJoin},
		"PART":{ 1, HandlePart},
		"PRIVMSG":{ 2, HandlePrivMsg},
		"OPER":{ 0, HandleOper},
		"MODE":{ 0, HandleMode},
	}

)

func
HandleNick(a HndlArg) error {
	newNick := a.msg.args[1]

	_, nickExists := srv.users[newNick]
	if nickExists {
		log.Printf("Nick '%s' is already taken\n", newNick)
		return message.Send(a.usr, Message{srv.host,
			[]string{FmtRplNum(ERR_NICKNAMEINUSE), "Nickname is already in use."},
		})
	}

	// Delete old user.
	if a.usr.nick != "" {
		delete(srv.users, a.usr.nick)
	}

	log.Printf("Set nick of '%s' to '%s'\n", a.usr.nick, newNick)
	a.usr.nick = newNick
	srv.users[newNick] = a.usr
	return nil
}

func
HandlePass(arg HndlArg) error {
	return nil
}

func
HandleSquit(arg HndlArg) error {
	return nil
}

func
HandleUser(a HndlArg) error {
	user, mode, _, info :=
		a.msg.args[1], a.msg.args[2],
		a.msg.args[3], a.msg.args[4]
	
	a.usr.user = user
	if v, err := strconv.Atoi(mode) ; err != nil {
		a.usr.mode = 0
	} else {
		a.usr.mode = v
	}
	a.usr.info = info
	
	return nil
}

func
HandleJoin(a HndlArg) error {
	chanStr := a.msg.args[1]
	chanNames := strings.Split(chanStr, ChanNameDelim)
	for _, v := range chanNames {
		// Skip channel names without prefixes.
		if format.HasAnyOfPrefixes(v, ChanNamePrefixes) == "" {
			continue	
		}

		ch, ok := srv.chans[v]
		// Create new channel if does not exist.
		if !ok {
			srv.chans[v] = &Channel{ make(map[string]*User)}
			ch = srv.chans[v]
		}

		ch.users[a.usr.nick] = a.usr
	}
	return nil
}

func
HandlePart(arg HndlArg) error {
	return nil
}

func
HandlePrivMsg(a HndlArg) error {
	var recvs []*User
	alltos := a.msg.args[1]
	msgstr := a.msg.args[2]
	names := strings.Split(alltos, ",")

	// Getting list of receivers.
	for _, to := range names {
		pref := format.HasAnyOfPrefixes(to, ChanNamePrefixes)
		if pref != "" { // For channels.
			ch, ok := srv.chans[to]
			if ok {
				for  _, v := range ch.users {
					recvs = append(recvs, v)
				}
			}
		} else { // For exact user.
			usr, ok := srv.users[to]
			if !ok {
				continue
			}
			recvs = append(recvs, usr)
		}
	}
	
	if len(recvs) == 0 {
		return SendMessage(a.usr,
			Message{
				a.msg.src,
				[]string{
					FmtRplNum(ERR_NORECIPIENT),
					fmt.Sprintf("No recipient given (%s)", a.msg.args[0]),
				},
			},
		)
	}

	// Sending to every of them.
	for _, u := range recvs {
		log.Printf("Sending private message to '%s'\n", a.usr.nick)
		err := SendMessage(
			u,
			Message{u.FullSrc(), []string{a.msg.args[0], a.usr.nick, msgstr}})
		if err == CIC {
			CleanUpUser(u)
		}
	}
	return nil
}
func
HandleOper(arg HndlArg) error {
	return nil
}

func
HandleMode(arg HndlArg) error {
	return nil
}


func
HandleMessage(usr *User, msg Message) error {
	cmd, ok := ClientCommands[msg.args[0]]
	if !ok {
		return errors.New("No such command")
	}

	if cmd.nargs > len(msg.args) - 1 {
		log.Printf("Not enough arguments for '%s'\n", msg.args[0])
		return nil
	}

	log.Printf("Handling '%s'", msg.args[0])
	err := cmd.hndl(HndlArg{usr, msg})
	if err != nil {
		return err
	}

	return nil
}

func
CleanUpUser(usr *User){
	/* We do not care if it is closed already. */
	usr.conn.Close()

	if usr.nick != "" {
		delete(srv.users, usr.nick)
	}
	log.Printf("Done cleaning '%s' user\n", usr.nick)
}

func
HandleConn(conn net.Conn) {
	/* New user. Not in server lists yet though. */
	usr := User{conn: conn}
	for {
		msg, err := ReadMsg(conn)
		/* If connection was closed in other thread 
			then also CIC will be returned*/
		if err == CIC {
			CleanUpUser(&usr)
			return;
		}
		fmt.Printf("'%s' %v\n", msg.src, msg.args)

		/* Handling includes writing replies,
			so CIC is checked when writing. */
		err = HandleMessage(&usr, msg)
		if err == CIC {
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

	srv = Server{host: host, port: port,
		ln: ln,
		users: make(map[string]*User),
		chans: make(map[string]*Channel), }

	for {
		conn, err := srv.ln.Accept()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(conn.RemoteAddr())
		go HandleConn(conn);
	}
}
