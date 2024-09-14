package main

import (
	"context"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os/signal"
	"path"
	"sync"
	"syscall"
)

func main() {
	//Create shutdown context to ensure graceful shutdown
	shutdownCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	//Parse cli arguments to application
	numOfConnections := flag.Int("n", 0, "connections to open")
	flag.Parse()
	parameters := flag.Args()
	serverURL := parameters[0]

	//Create the requested number of connections and wait for the app to exit
	var wg sync.WaitGroup
	for i := 0; i < *numOfConnections; i++ {
		wg.Add(1)
		go connectToWS(serverURL, &wg, i, shutdownCtx)
	}
	wg.Wait()
	defer cancel()
}

func connectToWS(serverURL string, wg *sync.WaitGroup, id int, shutdownCtx context.Context) {
	defer wg.Done()
	connURL := url.URL{
		Scheme: "ws",
		Host:   serverURL,
		Path:   path.Join("goapp", "ws"),
	}
	conn, _, err := websocket.DefaultDialer.Dial(connURL.String(), nil)
	if err != nil {
		log.Printf("could not connect to server because of %v\n", err)
		return
	}
	defer func(conn *websocket.Conn) {
		e := conn.Close()
		if e != nil {
			log.Printf("failed to close connection %d, because of %v", id, e)
		}
	}(conn)

	for {
		select {
		case <-shutdownCtx.Done():
			return
		default:
			_, incomingMsg, e := conn.ReadMessage()
			if e != nil {
				log.Printf("could not read from websocket because of %v\n", err)
				return
			}
			log.Printf("[conn %d] %s\n", id, string(incomingMsg))
		}
	}
}
