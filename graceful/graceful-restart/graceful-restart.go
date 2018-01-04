package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	w.Write([]byte("Hello, graceful restart!\n"))
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

	if isMaster() {
		log.Printf("master pid: %d\n", os.Getpid())
		laddr, _ := net.ResolveTCPAddr("tcp", "localhost:8888")
		l, _ := net.ListenTCP("tcp", laddr)
		supervise(l)
	}

	log.Printf("worker pid: %d\n", os.Getpid())
	fdStr := os.Getenv("__MASTER__")
	fd, _ := strconv.Atoi(fdStr)
	file := os.NewFile(uintptr(fd), "listen socket")
	defer file.Close()
	l, _ := net.FileListener(file)
	waitSignal(l)
	err := srv.Serve(l)
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

func isMaster() bool {
	return os.Getenv("__MASTER__") == ""
}

func supervise(l *net.TCPListener) {
	p, _ := forkExec(l)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR2)
	for {
		switch sig := <-c; sig {
		case syscall.SIGUSR2:
			new, _ := forkExec(l)
			p.Signal(syscall.SIGINT)
			p.Wait()
			p = new
		}
	}
}

func forkExec(l *net.TCPListener) (*os.Process, error) {
	progName, err := exec.LookPath(os.Args[0])
	if err != nil {
		return nil, err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	f, err := l.File()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	files := []*os.File{os.Stdin, os.Stdout, os.Stderr, f}
	fdEnv := fmt.Sprintf("%s=%d", "__MASTER__", len(files)-1)
	return os.StartProcess(progName, os.Args, &os.ProcAttr{
		Dir:   pwd,
		Env:   append(os.Environ(), fdEnv),
		Files: files,
	})
}
