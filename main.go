package main

import(
	"os"
	"log"
	"net"
	"fmt"
	"io"
	"errors"
	"bufio"
)

type User struct {
	Login string
}

var(
	MsgDelim = []byte("\r\n")
	MaxMsgLen = 512
	AddrStr = "localhost:6667"
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
ReadRawMsg(conn net.Conn) []byte {
	msg, _ := ReadTrunc(bufio.NewReader(conn), MaxMsgLen, MsgDelim)
	return msg
}

/*func
ReadMsg(conn net.Conn) []string {
	
}*/

func
HandleConn(conn net.Conn) {
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
