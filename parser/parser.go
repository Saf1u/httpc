package parser

import (
	"bufio"
	"errors"
	"fmt"
	"net445/connection"
	"strconv"
	"strings"
)

var headers = make(map[string]string)
var opts = make(map[string]bool)
var url string
var method string
var file string
var body string

func parseUrl(s string) string {
	r := strings.Split(s, "'http://")

	if len(r) == 2 {
		url = r[1]
		return "url"
	}
	return "nop"
}
func parseHeader(s string) string {
	list := strings.Split(s, ":")
	if len(list) != 2 || s[0:1] == "'" {
		return "nop"
	} else {
		headers[list[0]] = list[1]
		return "key:val"
	}

}

func parseBody(s string) string {
	if strings.HasPrefix(s, "'") {

		return "str"

	}
	if strings.HasPrefix(s, "safwan/") {
		file = s
		return "file"
	}
	return "nop"
}

func parseText(s string, ln *bufio.Scanner) string {
	if s[len(s)-1:] == "'" {
		body = s
		return "str"
	}
	text := make([]string, 0)
	text = append(text, s)

	for ln.Scan() {
		if ln.Text() == "'" {
			text = append(text, ln.Text())
			body = strings.Join(text, "")
			fmt.Println(body)
			return "str"
		} else {
			text = append(text, ln.Text())
		}
	}

	return ""
}

func parseflagormethod(s string) string {
	switch {
	case s == "-f":
		opts[s] = true
		return "f"
	case s == "-v":
		opts[s] = true
		return "v"
	case s == "-d":
		opts[s] = true
		return "d"
	case s == "-h":
		return "h"
	case s == "get":
		method = "GET"
		return "get"
	case s == "post":
		method = "POST"
		return "post"
	}
	return "nop"
}

func Parse(ln *bufio.Scanner) (*connection.Request, error) {
	symbolTable := make(map[string]map[string]string)
	for i := 0; i <= 25; i++ {
		num := strconv.Itoa(i)
		str := "q" + num
		symbolTable[str] = make(map[string]string)

	}
	symbolTable["qi"] = make(map[string]string)
	symbolTable["qi"]["httpc"] = "q0"

	symbolTable["q0"]["sp"] = "q1"
	symbolTable["q1"]["sp"] = "q1"
	symbolTable["q1"]["get"] = "q2"
	symbolTable["q2"]["sp"] = "q3"
	symbolTable["q3"]["sp"] = "q3"
	symbolTable["q3"]["sp"] = "q3"
	symbolTable["q3"]["h"] = "q6"
	symbolTable["q10"]["sp"] = "q10"
	symbolTable["q3"]["v"] = "q4"
	symbolTable["q4"]["sp"] = "q5"
	symbolTable["q5"]["sp"] = "q5"
	symbolTable["q5"]["h"] = "q6"

	symbolTable["q5"]["url"] = "q10"
	symbolTable["q6"]["sp"] = "q7"
	symbolTable["q7"]["sp"] = "q7"
	symbolTable["q7"]["key:val"] = "q8"
	symbolTable["q8"]["sp"] = "q9"
	symbolTable["q9"]["sp"] = "q9"
	symbolTable["q9"]["h"] = "q6"
	symbolTable["q9"]["url"] = "q10"
	symbolTable["q1"]["post"] = "q11"
	symbolTable["q11"]["sp"] = "q12"
	symbolTable["q12"]["sp"] = "q12"
	symbolTable["q12"]["v"] = "q13"
	symbolTable["q12"]["h"] = "q22"
	symbolTable["q22"]["sp"] = "q23"
	symbolTable["q23"]["sp"] = "q23"
	symbolTable["q23"]["key:val"] = "q24"

	symbolTable["q24"]["sp"] = "q25"
	symbolTable["q25"]["sp"] = "q25"
	symbolTable["q25"]["h"] = "q22"
	symbolTable["q25"]["d"] = "q15"
	symbolTable["q25"]["f"] = "q16"
	symbolTable["q13"]["sp"] = "q14"
	symbolTable["q14"]["sp"] = "q14"
	symbolTable["q14"]["d"] = "q15"
	symbolTable["q14"]["f"] = "q16"
	symbolTable["q14"]["sp"] = "q14"
	symbolTable["q14"]["h"] = "q22"

	symbolTable["q16"]["sp"] = "q18"
	symbolTable["q18"]["sp"] = "q18"
	symbolTable["q17"]["sp"] = "q17"
	symbolTable["q15"]["sp"] = "q17"
	symbolTable["q17"]["str"] = "q19"
	symbolTable["q18"]["file"] = "q20"
	symbolTable["q19"]["sp"] = "q21"
	symbolTable["q20"]["sp"] = "q21"
	symbolTable["q21"]["sp"] = "q21"
	symbolTable["q21"]["url"] = "q10"

	state := "qi"
	transition := ""
	ln.Split(bufio.ScanRunes)
	words := make([]string, 0)
	for ln.Scan() {

		word := ln.Text()
		if word != " " {
			words = append(words, word)
		} else {
			finalword := strings.Join(words, "")
			//fmt.Println(finalword)
			switch {
			case parseUrl(finalword) == "url":
				transition = "url"
			case parseHeader(finalword) == "key:val":
				transition = parseHeader(finalword)
			case parseflagormethod(finalword) != "nop":
				transition = parseflagormethod(finalword)

			case len(finalword) == 0:
				transition = "sp"
			case parseBody(finalword) != "nop":
				if parseBody(finalword) == "str" {
					transition = parseText(finalword, ln)
				} else {
					transition = "file"
				}

			case finalword == "httpc":
				transition = "httpc"
			}
			state = symbolTable[state][transition]
			if transition != "sp" {
				state = symbolTable[state]["sp"]
			}

			if state == "" {
				return nil, errors.New("bad command!")
			} else {
				//fmt.Println(state)
			}
			words = make([]string, 0)
		}
	}
	if state != "q10" {
		return nil, errors.New("bas command!")
	}
	//fmt.Println(state)
	host, resouce := getHostAndresource(url)

	req := &connection.Request{

		Method:          method,
		Host:            host,
		Address:         host,
		Port:            "80",
		ResourcePath:    resouce,
		ProtocolVersion: "HTTP/1.0",
		Headers:         headers,
		File:            file,
		Body:            body,
		Opts:            opts,
	}

	if body != "" {
		headers["Content-Length"] = strconv.Itoa(len(body) - 2)
	}

	return req, nil
}

func getHostAndresource(s string) (string, string) {
	s = strings.TrimSuffix(s, "'")
	ind := strings.Index(s, "/")
	host := s[0:ind]
	resouce := s[ind:]
	return host, resouce
}
