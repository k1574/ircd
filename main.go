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

type User struct {
	Login string
}

var(
	AddrStr = "localhost:6667"
	SrcPrefix = ":"
	ArgSep = " "
	LongArgSep =  ArgSep+":"
	MsgDelim = []byte("\r\n")
	MaxMsgLen = 512
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
ReadMsg(conn net.Conn) (string, []string, string, error) {
	buf, err := ReadRawMsg(conn)
	if err != nil {
		return "", nil, "", err
	}

	s := string(buf)

	src := ""
	if strings.HasPrefix(s, SrcPrefix) {
		s = s[len(SrcPrefix):]
		strs := strings.SplitN(s, ArgSep, 2)
		src, s = strs[0], strs[1]
	}

	args, lngarg := SplitTilSep(s, ArgSep, LongArgSep)
	return src, args, lngarg, nil
}

func
HandleConn(conn net.Conn) {
	for {
		pref, args, lngArgs, err := ReadMsg(conn)
		if err != nil {
			return;
		}
		fmt.Printf("'%s' %v '%s'\n", pref, args, lngArgs)
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
		go HandleConn(conn);
	}
}
