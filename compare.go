package main

import (
	"context"

	gorilla "github.com/gorilla/websocket"
	nhooyr "nhooyr.io/websocket"
)

// pseudo code
func Nhooyr() {
	var conn *nhooyr.Conn
	conn.Close(nhooyr.StatusNormalClosure, "") // send "close" message to peer

	var ctx context.Context
	nc := nhooyr.NetConn(ctx, conn, nhooyr.MessageText)
	nc.Close() // send "close" message to peer
}

// pseudo code
func Gorilla() {
	var conn *gorilla.Conn
	conn.Close() // closes the underlying network connection

	h := conn.CloseHandler()
	h(gorilla.CloseNormalClosure, "") // send "close" message to peer
}
