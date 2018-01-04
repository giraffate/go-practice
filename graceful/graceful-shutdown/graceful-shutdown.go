package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	w.Write([]byte("Hello, graceful shutdown!\n"))
}

func main() {
	http.HandleFunc("/", indexHandler)

	srv := &http.Server{Addr: ":8888", Handler: http.DefaultServeMux}

	var wg sync.WaitGroup
	srv.ConnState = func(conn net.Conn, state http.ConnState) {
		switch state {
		case http.StateActive:
			log.Println("StateActive!!!")
			wg.Add(1)
		case http.StateIdle:
			log.Println("StateIdle!!!")
			wg.Done()
		}
	}
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	waitSignal(l)
	err = srv.Serve(l)
	wg.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func waitSignal(l net.Listener) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go func() {
		sig := <-c
		switch sig {
		case syscall.SIGINT:
			signal.Stop(c)
			l.Close()
		}
	}()
}
