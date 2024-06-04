package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWSOrderbook(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client to orderbook feed:", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("orderbook data-> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 2)
	}
}

// this is for the chat
func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from clienit", ws.RemoteAddr())
	// each time the client is going to connect to the handler and we are going
	// to keep track of that connection
	// we need to write it
	s.conns[ws] = true

	// we are going to listen to messages to they send
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
				// conneciton on the other side has closed itself

			}
			fmt.Println("read error", err)
			continue
		}
		msg := buf[:n]
		//fmt.Println(string(msg))
		// we are going to make a chat real quick
		// we can reply with
		//ws.Write([]byte("thank you for the msg !!!"))
		s.broadcast(msg)
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error", err)
			}
		}(ws)
	}
}

func main() {

	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/orderbookfeed", websocket.Handler(server.handleWSOrderbook))
	http.ListenAndServe(":3000", nil)
}
