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
)

/* General type for interaction between client and server. */
type Message struct {
	src string
	args []string
	lngarg string
}

/* Commands from client. */
type Commands map[string] struct {
	nargs int
	hndl func(arg HndlArg)
}

/* Argument for handlers. */
type HndlArg struct {
	usr *User
	msg Message
}

/* Keep users and channels. */
type Server struct {
	users map[string]*User
	chans []*Channel
}

/* Stick connection and nick together. */
type User struct {
	conn net.Conn
	nick string
}

type Channel struct {
	users []User
}

var(
	srv Server
	AddrStr = "localhost:6667"

	MsgSrcPref = ":"
	MsgArgSep = " "
	MsgLongArgSep = MsgArgSep+":"
	MsgDelim = []byte("\r\n")
	MsgEmptyArg = "*"

	UserNameCantHave= []string{" "}
	ChanNamePrefixes = []string{"#", "&"}
	ChanNameCantHave = []string{" ", ",", string([]byte{7}) }

	ClientCommands = Commands{
		"NICK":{ 1, HandleNick},
		"PASS":{ 0, HandlePass},
		"SQUIT":{ 0, HandleSquit},
		"USER":{ 3, HandleUser},

		"JOIN":{ 1, HandleJoin},
		"PART":{ 1, HandlePart},
		"PRIVMSG":{ 1, HandlePrivMsg},
		"OPER":{ 0, HandleOper},
		"MODE":{ 0, HandleMode},
	}

)

const(
	RPL_WELCOME = 1
	RPL_YOURHOST = 2
	RPL_CREATED = 3
	RPL_MYINFO = 4
	RPL_BOUNCE = 5
	RPL_NONE = 300
	RPL_AWAY = 301
	RPL_USERHOST = 302
	RPL_ISON = 303
	RPL_UNAWAY = 305
	RPL_NOAWAY = 306
	RPL_WHOISUSER = 311
	RPL_WHOISSERVER = 312
	RPL_WHOISOPERATOR = 313
	RPL_WHOWASUSER = 314
	RPL_ENDOFWHO = 315
	RPL_WHOISIDLE = 317
	RPL_ENDOFWHOIS = 318
	RPL_WHOISCHANNELS = 319
	RPL_LIST = 322
	RPL_LISTEND = 323
	RPL_CHANNELMODEIS = 324
	RPL_UNIQOPIS = 325
	RPL_NOTOPIC = 331
	RPL_TOPIC = 332
	RPL_INVITING = 341
	RPL_SUMMONING = 342
	RPL_INVITELIST = 346
	RPL_ENDOFINVITELIST = 347
	RPL_EXCEPTLIST = 348
	RPL_ENDOFEXCEPTLIST = 349
	RPL_VERSION = 351
	RPL_WHOREPLY = 352
	RPL_NAMREPLY = 353
	RPL_LINKS = 364
	RPL_ENDOFLINKS = 365
	RPL_ENDOF_NAMES = 366
	RPL_BANLIST = 367
	RPL_ENDOFBANLIST = 368
	RPL_ENDOFWHOWAS = 369
	RPL_INFO = 371
	RPL_MOTD = 372
	RPL_ENDOFINFO = 374
	RPL_MOTDSTART = 375
	ERR_NICKNAMEINUSE = 433
)

const(
	MaxMsgLen = 512
	MaxChanNameLen = 200
	MaxClientNickLen = 9
)

func
HandleNick(a HndlArg) {
	newNick := a.msg.args[0]

	_, ok := srv.users[newNick]
	if ok {
		return
	}

	// Delete old user.
	if a.usr.nick != "" {
		delete(srv.users, a.usr.nick)
	}

	//fmt.Println("it fucking worked")
	a.usr.nick = newNick
	srv.users[newNick] = a.usr
}

func
(usr *User)SendMessage(msg Message) error {
	n, err := fmt.Fprint(usr.conn, MessageToRaw(msg))

	if err != nil {
		return err
	} else if n == 0 {
		return errors.New("Connection is closed")
	}

	return nil
}

func
MessageToRaw(msg Message) []byte {
	str := MsgSrcPref + msg.src +
		strings.Join(msg.args, MsgArgSep) +
		MsgLongArgSep + msg.lngarg +
		string(MsgDelim)
	return []byte(str)
}

func
HandlePass(arg HndlArg){
}

func
HandleSquit(arg HndlArg){
}

func
HandleUser(arg HndlArg){
}

func
HandleJoin(arg HndlArg){
}

func
HandlePart(arg HndlArg){
}

func
HandlePrivMsg(arg HndlArg){
}

func
HandleOper(arg HndlArg){
}

func
HandleMode(arg HndlArg){
}

/* Returns truncated bytes by size or delimiter. */
func
ReadTrunc(rd *bufio.Reader, siz int, delim []byte) ([]byte, error) {
	dlen := len(delim)	
	if dlen <= 0 {
		return nil, errors.New("delimiter length cannot be 0 or less")
	}
	var(
		ret []byte
		peakLen = siz - dlen 
		b byte
		buf []byte
	)

	buf = make([]byte, 1)
	
	i := 0
	j := 0
	for ;  i < peakLen ; i++ {
		n, err := rd.Read(buf)

		if n == 0 {
			return nil, errors.New("read 0 bytes, so leaving")
		} else if err == io.EOF {
			break;
		} else if err != nil {
			return nil, err
		}

		b = buf[0]
		ret = append(ret, b)
		if b == delim[j] {
			if j == dlen - 1 { break }
			j++
		} else {
			j = 0
		}
	}

	if i == peakLen {
		ret = append(ret, delim...)
	}
	
	
	return ret, nil
}

func
SplitTilSep(s, sep, endsep string) ([]string, string) {
	n := strings.LastIndex(s, endsep)
	var arg, str string
	if n != -1 {
		arg = s[:n]
		str = s[n+len(endsep):]
	} else {
		arg = s
		str = ""
	}

	return strings.Split(arg, sep), str
}

func
ReadRawMsg(conn net.Conn) ([]byte, error) {
	msg, err := ReadTrunc(bufio.NewReader(conn), MaxMsgLen, MsgDelim)
	return msg, err
}

func
ReadMsg(conn net.Conn) (Message, error) {
	buf, err := ReadRawMsg(conn)
	if err != nil {
		return Message{"", nil, ""}, err
	}

	s := string(buf)

	src := ""
	if strings.HasPrefix(s, MsgSrcPref) {
		s = s[len(MsgSrcPref):]
		strs := strings.SplitN(s, MsgArgSep, 2)
		src, s = strs[0], strs[1]
	}

	s = s[:len(s)-len(MsgDelim)]
	args, lngarg := SplitTilSep(s, MsgArgSep, MsgLongArgSep)
	return Message{src, args, lngarg}, nil
}

func
HandleMessage(usr *User, msg Message) error {
	if(len(msg.args) < 1){
		return nil
	}
	
	cmd, ok := ClientCommands[msg.args[0]]
	if !ok {
		return errors.New("No such command")
	}

	if cmd.nargs != len(msg.args) - 1 {
		//FMT.Println("fuck you")
		return nil
	}

	cmd.hndl(HndlArg{usr, msg})

	return nil
}

func
HandleConn(conn net.Conn) {
	usr := User{conn, ""}
	for {
		msg, err := ReadMsg(conn)
		if err != nil {
			return;
		}
		//fmt.Printf("'%s' %v '%s'\n", msg.src, msg.args, msg.lngarg)
		HandleMessage(&usr, msg)
	}
}

func
main() {
	ln, err := net.Listen("tcp", AddrStr)
	if err != nil {
		log.Fatal(err)
	}
	srv.users = make(map[string]*User)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(conn.RemoteAddr())
		go HandleConn(conn);
	}
}
