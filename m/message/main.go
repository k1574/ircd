package message

/* General type for interaction between client and server. */
type Message struct {
	src string
	args []string
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
	lastIdx := len(msg.args) - 1
	str := ""

	/* Handle long last arguments. */
	if strings.Contains(msg.args[len(msg.args)-1], MsgArgSep) {
		hasLongArg = true
		if lastIdx != 0 {
			lastIdx--
		}
	}

	if msg.src != "" {
		str += SrcPref + msg.src
	}

	if lastIdx > 0 {
		str += ArgSep
		str += strings.Join(msg.args[:lastIdx], MsgArgSep)
	}

	if hasLongArg {
		str += MsgLongArgSep + msg.args[len(msg.args)-1]
	}

	str += string(MsgDelim)

	return []byte(str)
}

func
ReadRaw(conn net.Conn) ([]byte, error) {
	msg, err := ReadTrunc(bufio.NewReader(conn), MaxMsgLen, MsgDelim)
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
	if strings.HasPrefix(s, MsgSrcPref) {
		s = s[len(MsgSrcPref):]
		strs := strings.SplitN(s, MsgArgSep, 2)
		src, s = strs[0], strs[1]
	}

	s = s[:len(s)-len(MsgDelim)]
	args, lngarg := SplitTilSep(s, MsgArgSep, MsgLongArgSep)
	if lngarg != "" {
		args = append(args, lngarg)
	}
	return Message{src, args}, nil
}


func
Send(usr *user.User, msg Message) error {
	n, err := fmt.Fprint(usr.conn, string(MessageToRaw(msg)))

	if err != nil {
		return err
	} else if n == 0 {
		return CIC
	}

	return nil
}

