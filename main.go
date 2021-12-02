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
	hndl func(arg HndlArg) error
}

/* Argument for handlers. */
type HndlArg struct {
	usr *User
	msg Message
}

/* Keep users and channels. */
type Server struct {
	host string
	port int
	ln net.Listener
	users map[string]*User
	chans map[string]*Channel
}

/* Stick connection and nick together. */
type User struct {
	conn net.Conn
	nick, user, host, info string
}

type Channel struct {
	users []*User
}

var(
	srv Server
	// Conection is closed error
	CIC = errors.New("connection is closed")

	MsgSrcPref = ":"
	MsgArgSep = " "
	MsgLongArgSep = MsgArgSep+":"
	MsgDelim = []byte("\r\n")
	MsgEmptyArg = "*"

	ChanNamePrefixes = []string{"#", "&"}
	ChanNameCantHave = []string{",", string([]byte{7}) }

	UserNameCantHave= append(ChanNamePrefixes, []string{}...)

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
	RPL_TRACELINK = 200
	RPL_TRACECONNECTING = 201
	RPL_TRACEHANDSHAKE = 202
	RPL_TRACEUNKNOWN = 203
	RPL_TRACEOPERATOR = 204
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
	ERR_NOADMININFO = 423
	ERR_FILEERROR = 424
	ERR_NO_NICKNAMEGIVEN = 431
	ERR_ERRONEUSNICKNAME = 432
	ERR_NICKNAMEINUSE = 433
	ERR_NICKCOLLISION = 436
	ERR_UNAVAILSOURCE = 437
	ERR_USERNOTINCHANNEL = 441
	ERR_NOTONCHANNEL = 442
	ERR_USERONCHANNEL = 443
)

const(
	MaxMsgLen = 512
	MaxChanNameLen = 200
	MaxClientNickLen = 9
)

func
HandleNick(a HndlArg) error {
	newNick := a.msg.args[1]

	_, ok := srv.users[newNick]
	if ok {
		log.Printf("Nick '%s' is already taken\n", newNick)
		return nil
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
FmtRplNum(num int) string {
	return fmt.Sprintf("%03d", num)
}

func
(u *User)FullSrc() string {
	return u.nick+"!"+u.user+"@"+u.host
}

func
SendMessageToUser(usr *User, msg Message) error {
	n, err := fmt.Fprint(usr.conn, string(MessageToRaw(msg)))

	if err != nil {
		return err
	} else if n == 0 {
		return CIC
	}

	return nil
}

func
MessageToRaw(msg Message) []byte {
	str := ""
	if msg.src != "" {
		str += MsgSrcPref + msg.src
	}

	str += strings.Join(msg.args, MsgArgSep)

	if msg.lngarg != "" {
		str += MsgLongArgSep + msg.lngarg
	}

	str += string(MsgDelim)

	return []byte(str)
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
HandleUser(arg HndlArg) error {
	return nil
}

func
HandleJoin(arg HndlArg) error {
	return nil
}

func
HandlePart(arg HndlArg) error {
	return nil
}

func
HasAnyOfPrefixes(s string, prefs []string) string {
	for _, v := range prefs {
		if strings.HasPrefix(s, v) {
			return v
		}
	}
	return ""
}

func
HandlePrivMsg(a HndlArg) error {
	to := a.msg.args[1]
	msgstr := a.msg.lngarg
	var recvs []*User

	// Getting list of receivers.
	pref := HasAnyOfPrefixes(to, ChanNamePrefixes)
	if pref != "" { // For channels.
		ch, ok := srv.chans[to]
		if !ok {
			return nil
		}
		recvs = ch.users
	} else { // For exact user.
		recvs = []*User{srv.users[to]}
	}

	// Sending to every of them.
	for _, u := range recvs {
		log.Printf("Sending private message to '%s'\n", a.usr.nick)
		err := SendMessageToUser(
			u,
			Message{"", []string{a.msg.args[0], a.usr.nick}, msgstr})
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
			return nil, CIC
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
	usr := User{conn, "", "", srv.host, "" }
	for {
		msg, err := ReadMsg(conn)
		/* If connection was closed in other thread 
			then also CIC will be returned*/
		if err == CIC {
			CleanUpUser(&usr)
			return;
		}
		//fmt.Printf("'%s' %v '%s'\n", msg.src, msg.args, msg.lngarg)

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
