package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net445/connection"
	"net445/parser"
	"strings"
)

type response struct {
	statusLine string //contains protocolMethod SP respcode SP phrase associated with code
	headers    map[string]string
	//resp Headers
	// Location
	// Server
	// WWW-Authenthicate

	//entity headers
	//                     Allow                    ; Section 10.1
	//                   | Content-Encoding         ; Section 10.3
	//                   | Content-Length           ; Section 10.4
	//                   | Content-Type             ; Section 10.5
	//                   | Expires                  ; Section 10.7
	//                   | Last-Modified            ; Section 10.10
	//                   | extension-header

}

func main() {
	rd := bufio.NewScanner(strings.NewReader("httpc post -h Content-Type:application/json -d '{\"Assignment\": 1}' 'http://httpbin.org/post' "))
	//"httpc post -v -h User-Agent:Mozilla/5.0 -d 'goo'  'http://httpbin.org/get?course=networking&assignment=1' "))
	req, err := parser.Parse(rd)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(req)
		httpTemplate := connection.BuildHttpTemplate(req)
		addr := req.Host + ":" + req.Port
		con, err := net.Dial("tcp", addr)
		defer con.Close()
		if err != nil {
			log.Println(err)
			return
		}
		err = connection.Send(httpTemplate, con)
		if err != nil {
			log.Println(err)
			return
		}
		buf := make([]byte, 1024)
		err = connection.Receive(con, buf)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(buf))

	}

}
