package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net445/connection"
	"os"
	"strconv"
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


type headerFlag []string

func (h *headerFlag) String() string {
	return "none"
}

func (h *headerFlag) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func parse() *connection.Request {
	getCommmand := flag.NewFlagSet("get", flag.ExitOnError)
	postCommand := flag.NewFlagSet("post", flag.ExitOnError)
	var headers headerFlag
	verbose := false
	getCommmand.BoolVar(&verbose, "v", false, "verbose option")
	postCommand.BoolVar(&verbose, "v", false, "verbose option")
	getCommmand.Var(&headers, "h", "header flags")
	postCommand.Var(&headers, "h", "header flags")
	bodyString := ""
	postCommand.StringVar(&bodyString, "d", "", "string body")
	file := ""
	postCommand.StringVar(&file, "f", "", "file option")

	if len(os.Args) < 2 {
		fmt.Println("not recognized")
		fmt.Println("see help for usage")
		os.Exit(0)
	}
	switch os.Args[1] {
	case "get":
		if len(os.Args) >= 3 {
			if os.Args[2] == "help" {
				fmt.Println(get)
			} else {
				getCommmand.Parse(os.Args[2:])
			}
		} else {
			fmt.Println("not recognized")
			fmt.Println("see help for usage")
			os.Exit(0)
		}
	case "post":
		if len(os.Args) >= 3 {
			if os.Args[2] == "help" {
				fmt.Println(post)
			} else {
				postCommand.Parse(os.Args[2:])
			}
		} else {
			fmt.Println("not recognized")
			fmt.Println("see help for usage")
			os.Exit(0)
		}
	case "help":
		fmt.Println(help)
		os.Exit(0)
	default:
		fmt.Println("not recognized")
		fmt.Println("see help for usage")
		os.Exit(0)
	}

	if getCommmand.Parsed() {
		host, resource := getHostAndresource(os.Args[len(os.Args)-1])
		list := extractHeaders(headers)
		opts := make(map[string]bool)
		opts["verbose"] = verbose
		req := &connection.Request{

			Method:          "GET",
			Host:            host,
			Address:         host,
			Port:            "80",
			ResourcePath:    resource,
			ProtocolVersion: "HTTP/1.0",
			Headers:         list,
			Opts:            opts,
		}

		return req
	} else {
		if postCommand.Parsed() {
			if bodyString != "" && file != "" {
				fmt.Println("not recognized")
				fmt.Println("see help for usage")
				os.Exit(0)
			}
			url := os.Args[len(os.Args)-1]
			if len(url) < 8 {
				fmt.Println("not recognized")
				fmt.Println("see help for usage")
				os.Exit(0)
			}
			host, resource := getHostAndresource(url)
			list := extractHeaders(headers)
			opts := make(map[string]bool)
			opts["verbose"] = verbose
			req := &connection.Request{

				Method:          "POST",
				Host:            host,
				Address:         host,
				Port:            "80",
				ResourcePath:    resource,
				ProtocolVersion: "HTTP/1.0",
				Headers:         list,
				File:            file,
				Body:            bodyString,
				Opts:            opts,
			}
			if bodyString != "" {
				list["Content-Length"] = strconv.Itoa(len(req.Body))
			} else {
				fl, err := os.OpenFile(file, os.O_RDONLY, 755)
				if err != nil {
					fmt.Println("could not open file for reading!")
					os.Exit(0)
				}
				message, err := ioutil.ReadAll(fl)
				if err != nil {
					fmt.Println("could not read data!", err.Error())
					os.Exit(0)
				}
				req.File = string(message)
				list["Content-Length"] = strconv.Itoa(len(req.File))

			}

			return req

		}

	}
	return nil
}

func main() {
	req := parse()
	httpTemplate := connection.BuildHttpTemplate(req)

	addr := req.Host + ":" + req.Port
	fmt.Println(addr)
	con, err := net.Dial("tcp", addr)
	defer func() {
		if con != nil {
			con.Close()
		}
	}()
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
	stringRep := string(buf)
	fmt.Println(stringRep)
}

func extractHeaders(header headerFlag) map[string]string {
	list := make(map[string]string)
	for _, value := range header {
		str := strings.Split(value, ":")
		if len(str) < 2 {
			fmt.Println("not recognized")
			fmt.Println("see help for usage")
			os.Exit(0)
		}
		list[str[0]] = str[1]
	}
	return list
}

func getHostAndresource(s string) (string, string) {
	s = s[7:]
	ind := strings.Index(s, "/")
	host := s[0:ind]
	resouce := s[ind:]
	return host, resouce
}
