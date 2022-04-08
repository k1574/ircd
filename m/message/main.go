package message

import(
	"strings"
	"bufio"
	"net"
	"fmt"
	"ircd/m/input"
	"ircd/m/format"
	"ircd/m/user"
)

/* General type for interaction between client and server. */
type Message struct {
	Src string
	Args []string
}

var(
	SrcPref = ":"
	ArgSep = " "
	LngArgSep = ArgSep+":"
	Del = []byte("\r\n")
	EmptyArg = "*"
	MaxLen = 512
)

func
ToRaw(msg Message) []byte {
	hasLongArg := false
	/* Last index for commong, not long, arguments. */
	lastIdx := len(msg.Args) - 1
	str := ""

	/* Handle long last arguments. */
	if strings.Contains(msg.Args[lastIdx], ArgSep) {
		hasLongArg = true
		if lastIdx != 0 {
			lastIdx--
		}
	}

	if msg.Src != "" {
		str += SrcPref + msg.Src
	}

	if lastIdx > 0 {
		str += ArgSep
		str += strings.Join(msg.Args[:lastIdx], ArgSep)
	}

	if hasLongArg {
		str += LngArgSep + msg.Args[len(msg.Args)-1]
	}

	str += string(Del)

	return []byte(str)
}

func
ReadRaw(conn net.Conn) ([]byte, error) {
	msg, err := input.ReadTrunc(bufio.NewReader(conn), MaxLen, Del)
	return msg, err
}

func
Read(conn net.Conn) (Message, error) {
	buf, err := ReadRaw(conn)
	if err != nil {
		return Message{"", nil}, err
	}

	s := string(buf)

	src := ""
	if strings.HasPrefix(s, SrcPref) {
		s = s[len(SrcPref):]
		strs := strings.SplitN(s, ArgSep, 2)
		src, s = strs[0], strs[1]
	}

	s = s[:len(s)-len(Del)]
	args, lngarg := format.SplitTilSep(s, ArgSep, LngArgSep)
	if lngarg != "" {
		args = append(args, lngarg)
	}
	return Message{src, args}, nil
}


func
Send(usr *user.User, msg Message) error {
	n, err := fmt.Fprint(usr.Conn, string(ToRaw(msg)))

	if err != nil {
		return err
	} else if n == 0 {
		return input.CIC
	}

	return nil
}

