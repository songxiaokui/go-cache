package internal

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

/*
@Time    : 2021/3/4 20:31
@Author  : austsxk
@Email   : austsxk@163.com
@File    : server.go
@Software: GoLand
*/

type Server struct {
	Cache
}

// Listen tcp port reviver
func (s *Server) Listen() {
	client, err := net.Listen("tcp", ":12345")
	fmt.Println("tcp connect port=12345")
	if err != nil {
		log.Fatal("tcp listen error", err)
	}
	for {
		c, err := client.Accept()
		if err != nil {
			log.Fatal("tcp client accept error:", err)
		}
		// goroutine to deal with thing
		go s.process(c)
	}
}

// make function
func NewTcpClient(c *Cache) *Server {
	return &Server{*c}
}

// Server readKey
func (s *Server) readKey(r *bufio.Reader) (string, error) {
	klen, err := readLen(r)
	if err != nil {
		return "", err
	}
	k := make([]byte, klen)
	// read all date to k
	_, err = io.ReadFull(r, k)
	if err != nil {
		return "", err
	}
	return string(k), nil
}

// readLen get bufio.Reader len depend space
func readLen(r *bufio.Reader) (int, error) {
	tmp, err := r.ReadString(' ')
	if err != nil {
		return 0, err
	}
	l, err := strconv.Atoi(strings.TrimSpace(tmp))
	if err != nil {
		return 0, err
	}
	return l, nil
}

func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	// get key command key sp value sp
	klen, err := readLen(r)
	if err != nil {
		return "", nil, err
	}
	vlen, err := readLen(r)
	if err != nil {
		return "", nil, err
	}
	k := make([]byte, klen)
	_, err = io.ReadFull(r, k)
	if err != nil {
		return "", nil, err
	}
	v := make([]byte, vlen)
	_, err = io.ReadFull(r, v)
	if err != nil {
		return "", nil, err
	}
	return string(k), v, nil
}

// return response
func sendResponse(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error()
		tmp := fmt.Sprintf("-%d", len(errString)) + errString
		_, e := conn.Write([]byte(tmp))
		return e
	}
	vlen := fmt.Sprintf("%d", len(value))
	_, e := conn.Write(append([]byte(vlen), value...))
	return e
}

// get set del impl
func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	k, e := s.readKey(r)
	if e != nil {
		return e
	}
	v, err := s.Get(k)
	return sendResponse([]byte(v), err, conn)
}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	k, v, err := s.readKeyAndValue(r)
	if err != nil {
		return err
	}
	return sendResponse(v, s.Set(k, string(v)), conn)
}

func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	k, err := s.readKey(r)
	if err != nil {
		return err
	}
	err = s.Delete(k)
	return sendResponse(nil, err, conn)
}

// goroutine content
func (s *Server) process(conn net.Conn) {
	defer conn.Close()
	// from connect get data return a bufio.reader
	r := bufio.NewReader(conn)
	for {
		op, err := r.ReadByte()
		if err != nil {
			if err != io.EOF {
				log.Println("close connect due to connect no data")
				return
			}
			if op == 'S' {
				err = s.set(conn, r)
			} else if op == 'G' {
				err = s.get(conn, r)
			} else if op == 'D' {
				err = s.del(conn, r)
			} else {
				log.Println("close connect due to invalid operation", op)
				return
			}
			if err != nil {
				log.Println("close connect due to error:", err)
				return
			}
		}
	}
}
