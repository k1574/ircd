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
	pref string
	args []string
	lngarg string
}

/* Commands from client. */
type Command struct {
	name string
	nargs int
}

type Server struct {
	users []User
	chans []Channel
}

type User struct {
	login string
}

type Channel struct {
	users []User
}

var(
	srv Server
	AddrStr = "localhost:6667"

	MsgSrcPrefix = ":"
	MsgArgSep = " "
	MsgLongArgSep = MsgArgSep+":"
	MsgDelim = []byte("\r\n")
	MsgEmptyArg = "*"

	ChanNamePrefixes = []string{"#", "&"}
	ChanNameCantHave = []string{" ", ",", string([]byte{7}) }

	ClientCommands = []Command{
		{"NICK", 1},
		{"PASS", 0},
		{"SQUIT", 0},
		{"USER", 3},

		{"JOIN", 1},
		{"PART", 1},
		{"PRIVMSG", 1},
		{"OPER", 0},
		{"MODE", 0},
	}

)

const(
	RPL_WELCOME = 1
	RPL_YOURHOST = 2
	RPL_CREATED = 3
	RPL_MYINFO = 4
	RPL_BOUNCE = 5
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
	if strings.HasPrefix(s, MsgSrcPrefix) {
		s = s[len(MsgSrcPrefix):]
		strs := strings.SplitN(s, MsgArgSep, 2)
		src, s = strs[0], strs[1]
	}

	s = s[:len(s)-len(MsgDelim)]
	args, lngarg := SplitTilSep(s, MsgArgSep, MsgLongArgSep)
	return Message{src, args, lngarg}, nil
}

func
handleMessage(conn net.Conn, msg Message) {
	if(len(msg.args) < 1){
		return
	}

	switch(msg.args[0]){
	}
}

func
handleConn(conn net.Conn) {
	for {
		msg, err := ReadMsg(conn)
		if err != nil {
			return;
		}
		//fmt.Printf("'%s' %v '%s'\n", msg.pref, msg.args, msg.lngarg)
		handleMessage(conn, msg)
	}
}

func
main() {
	ln, err := net.Listen("tcp", AddrStr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(conn.RemoteAddr())
		go handleConn(conn);
	}
}
