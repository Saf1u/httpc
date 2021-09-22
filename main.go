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

var (
	get = `httpc help get
	usage: httpc get [-v] [-h key:value] URL
	Get executes a HTTP GET request for a given URL.
	 -v Prints the detail of the response such as protocol, status,
	and headers.
	 -h key:value Associates headers to HTTP Request with the format
	'key:value'.`

	post = `usage: httpc post [-v] [-h key:value] [-d inline-data] [-f file] URL
	Post executes a HTTP POST request for a given URL with inline data or from
	file.
	 -v Prints the detail of the response such as protocol, status,
	and headers.
	 -h key:value Associates headers to HTTP Request with the format
	'key:value'.
	 -d string Associates an inline data to the body HTTP POST request.
	 -f file Associates the content of a file to the body HTTP POST
	request.`
	help = `httpc help
    httpc is a curl-like application but supports HTTP protocol only.
    Usage:
    httpc command [arguments]
    The commands are:
 	get executes a HTTP GET request and prints the response.
 	post executes a HTTP POST request and prints the response.
 	help prints this screen.
	Use "httpc help [command]" for more information about a command.
	`
)

//"httpc post -v -h User-Agent:Mozilla/5.0 -d 'goo'  'http://httpbin.org/get?course=networking&assignment=1' "))
//rd := bufio.NewScanner(strings.NewReader("httpc post -h Content-Type:application/json -d '{"assignment": 1}' 'http://httpbin.org/post' "))

func main() {
	for {
		httpc()
	}

}

func httpc() {
	reader := bufio.NewReader(os.Stdin)
	query, err := reader.ReadString('\n')
	check := true
	for check {
		list := strings.TrimSpace(query)
		switch {
		case list == "httpc help":
			fmt.Println(help)
			query, _ = reader.ReadString('\n')
		case list == "httpc help get":
			fmt.Println(get)
			query, _ = reader.ReadString('\n')
		case list == "httpc help post":
			fmt.Println(post)
			query, _ = reader.ReadString('\n')
		default:
			check = false

		}

	}
	query = query[0 : len(query)-1]
	if err != nil {
		log.Println(err)
		return
	}
	query = query + " "
	rd := bufio.NewScanner(strings.NewReader(query))
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
