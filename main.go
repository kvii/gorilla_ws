package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

const u = "ws://localhost:9090"

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { return client(ctx) })
	eg.Go(func() error { return server(ctx) })

	err := eg.Wait()
	if err != nil {
		panic(err)
	}
}

func client(ctx context.Context) error {
	conn, resp, err := websocket.DefaultDialer.DialContext(ctx, u, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	defer conn.Close()

	wc, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	fmt.Fprint(wc, "ping")
	err = wc.Close()
	if err != nil {
		return err
	}

	_, r, err := conn.NextReader()
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		return err
	}

	f := conn.CloseHandler()
	err = f(websocket.CloseNormalClosure, "")
	if err != nil {
		return err
	}
	return nil
}

func server(ctx context.Context) error {
	var up websocket.Upgrader
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := up.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer conn.Close()

		err = handle(conn)
		if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	l, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return err
	}

	var s http.Server
	go s.Serve(l)
	<-ctx.Done()
	return s.Close()
}

func handle(conn *websocket.Conn) error {
	for {
		_, reader, err := conn.NextReader()
		if err != nil {
			return err
		}
		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			return err
		}

		wc, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return err
		}
		fmt.Fprint(wc, "pong")
		err = wc.Close()
		if err != nil {
			return err
		}
	}
}
