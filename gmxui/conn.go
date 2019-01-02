package main

import (
	"net"
	"encoding/json"
	"strings"
	"sync"
	"log"
)

// GMXConn is connection between gmx client and gmx process.
type GMXConn struct {
	addr string
	net.Conn
	*json.Decoder
	*json.Encoder
}

// FetchKeys returns all the registered keys from the process.
func (conn *GMXConn) FetchKeys() []string {
	// retrieve list of registered keys
	if err := conn.Encode([]string{"keys"}); err != nil {
		log.Fatalf("unable to send keys request to process: %v", err)
	}
	var result = make(map[string][]string)
	if err := conn.Decode(&result); err != nil {
		log.Fatalf("unable to decode keys response: %v", err)
	}
	keys, ok := result["keys"]
	if !ok {
		log.Fatalf("gmx server did not return a keys list")
	}
	return keys
}


// GetValues returns values of keys from the process.
func (conn *GMXConn) GetValues(keys []string) interface{} {
	// retrieve list of registered keys
	if err := conn.Encode(keys); err != nil {
		log.Fatalf("unable to send request to address: %v", err)
	}
	var result interface{}
	if err := conn.Decode(&result); err != nil {
		log.Fatalf("unable to decode response: %v", err)
	}
	
	return result
}

func dial(addr string) (*GMXConn, error) {
	var socketType string

	if strings.HasPrefix(addr, "/") {
		socketType = "unix"
	} else {
		socketType = "tcp"
	}

	c, err := net.Dial(socketType, addr)
	return &GMXConn{
		addr,
		c,
		json.NewDecoder(c),
		json.NewEncoder(c),
	}, err
}


type GMXConnPool struct {
	sync.Mutex
	conns map[string]*GMXConn	// tracks all connected GMX processes
} 

func NewGMXConnPool() *GMXConnPool {
	return &GMXConnPool{
		conns: make(map[string]*GMXConn),
	}
}

func (p *GMXConnPool) Push(addr string) error {
	p.Lock()
	defer p.Unlock()

	if _, exist := p.conns[addr]; exist {
		return nil
	}

	conn, err := dial(addr)
	if err != nil {
		return err
	}
	p.conns[addr] = conn
	
	return nil
}


func (p *GMXConnPool) Get(addr string) *GMXConn {
	return p.conns[addr]
}

func (p *GMXConnPool) HasAddr(addr string) bool {
	_, exist := p.conns[addr]
	return exist
}


