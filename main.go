package main

import(
	"os"
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
	Prefix = ":"
	LongArgSep = ":"
	ArgSep = " "
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
	var ret []byte
	var peakLen = siz - dlen 

	i := 0
	j := 0
	for ;  i < peakLen ; i++ {
		b, err := rd.ReadByte()

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

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
		str = s[n:]
	} else {
		arg = s
		str = ""
	}

	return strings.Split(arg, sep), str
}

func
ReadRawMsg(conn net.Conn) []byte {
	msg, _ := ReadTrunc(bufio.NewReader(conn), MaxMsgLen, MsgDelim)
	return msg
}

func
HandleConn(conn net.Conn) {
}

func
main() {
	rd := bufio.NewReader(os.Stdin)
	val, _ := ReadTrunc(rd, 20, MsgDelim)
	args, str := SplitTilSep(string(val), ArgSep, ArgSep + LongArgSep)
	fmt.Printf("%v %d '%s'", args, len(args), str)
	return
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
