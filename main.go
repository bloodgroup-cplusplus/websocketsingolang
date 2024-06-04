package main

import "golang.org/x/net/websocket"

type Server struct {
	conns map[*websocket.Conn]bool
}

func main() {

}
