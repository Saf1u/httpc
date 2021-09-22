package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net445/connection"
	"net445/parser"
	"os"
	"strings"
)

func main() {
	//rd := bufio.NewScanner(strings.NewReader("httpc post -h Content-Type:application/json -d '{\"assignment\": 1}' 'http://httpbin.org/post' "))
	reader := bufio.NewReader(os.Stdin)
	query, err := reader.ReadString('\n')
	query = query[0 : len(query)-1]
	if err != nil {
		log.Println(err)
		return
	}
	query = query + " "
	rd := bufio.NewScanner(strings.NewReader(query))
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
