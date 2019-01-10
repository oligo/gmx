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
		log.Printf("unable to send keys request to process: %v\n", err)
		return nil
	}

	var result = make(map[string][]string)
	if err := conn.Decode(&result); err != nil {
		log.Printf("unable to decode keys response: %v\n", err)
		return nil
	}

	keys, ok := result["keys"]
	if !ok {
		log.Printf("gmx server did not return a keys list\n")
	}
	
	return keys
}


// GetValues returns values of keys from the process.
func (conn *GMXConn) GetValues(keys []string) map[string]interface{} {
	// retrieve list of registered keys
	if err := conn.Encode(keys); err != nil {
		log.Fatalf("unable to send request to address: %v", err)
		return nil
	}

	var result map[string]interface{}
	if err := conn.Decode(&result); err != nil {
		log.Fatalf("unable to decode response: %v", err)
		return nil
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

func (p *GMXConnPool) Push(addr string) (*GMXConn, error) {
	p.Lock()
	defer p.Unlock()

	if conn, exist := p.conns[addr]; exist {
		return conn, nil
	}

	conn, err := dial(addr)
	if err != nil {
		return nil, err
	}
	p.conns[addr] = conn
	
	return conn, nil
}


func (p *GMXConnPool) Get(addr string) *GMXConn {
	return p.conns[addr]
}

func (p *GMXConnPool) HasAddr(addr string) bool {
	_, exist := p.conns[addr]
	return exist
}


