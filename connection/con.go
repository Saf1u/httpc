package connection

import (
	"errors"
	"net"
)

type Request struct {
	Method          string
	Host            string
	Address         string
	Port            string
	ResourcePath    string
	ProtocolVersion string
	Headers         map[string]string
	File            string
	Body            string
	DumpFile        string
	Opts            map[string]bool
}

func BuildHttpTemplate(req *Request) string {
	s := ""
	if req.Method == "GET" {
		s = buildGet(req)
	} else {
		s = buildPost(req)
	}
	return s
}

func Send(s string, connection net.Conn) error {
	_, err := connection.Write([]byte(s))
	if err != nil {
		return errors.New("could not write on tcp:" + err.Error())
	}
	return nil
}

func Receive(connection net.Conn, buf []byte) error {

	_, err := connection.Read(buf)
	if err != nil {
		return errors.New("could not read from tcp:" + err.Error())
	}
	return nil
}

func buildGet(req *Request) string {
	s := "GET " + req.ResourcePath + " " + req.ProtocolVersion + "\r\n" + "Host: " + req.Host + "\r\n"
	for key, val := range req.Headers {
		s = s + key + ":" + val + "\r\n"
	}
	s = s + "\r\n"
	return s
}

func buildPost(req *Request) string {
	body := ""
	if req.Body == "" {
		body = req.File
	} else {
		body = req.Body
	}
	s := "POST " + req.ResourcePath + " " + req.ProtocolVersion + "\r\n" + "Host: " + req.Host + "\r\n"
	for key, val := range req.Headers {
		s = s + key + ":" + val + "\r\n"
	}
	s = s + "\r\n" + body + "\r\n"

	return s
}
