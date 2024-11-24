package main

import (
	postpkg "app-go/internal/modules/post/services"
	"app-go/internal/servers/http"
	"app-go/internal/servers/http/handlers"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// initializing pre-requisites
	log.Println("starting server...")

	// NOTE: add modules initialization here
	// ...

	// initialize modukes
	postmodule := postpkg.New(nil, nil)

	// NOTE: add server initialization on the following
	// ...
	posthandler := handlers.NewPost(postmodule)

	httpserver := http.New(":4000", posthandler)

	// server starts
	httpserver.Start()

	waitForExitSignal()

	// shutdown starts
	log.Println("shutting down server...")

	// NOTE: add shutdown handler on the following
	// ...

	httpserver.Stop() // stopping http server

	log.Println("server stopped")
}

// waitForExitSignal is to awaits incoming interrupt signal
// sent to the service
func waitForExitSignal() os.Signal {
	log.Printf("awaiting exit signal...")
	ch := make(chan os.Signal, 4)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	)

	return <-ch
}
