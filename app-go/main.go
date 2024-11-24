package main

import (
	"app-go/internal/clients/postgres"
	filestorageS3pkg "app-go/internal/modules/file_storage/s3"
	postpkgrepo "app-go/internal/modules/post/repositories/postgres"
	postpkg "app-go/internal/modules/post/services"
	"app-go/internal/servers/http"
	"app-go/internal/servers/http/handlers"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	log.Println("starting server...")

	// NOTE: add base initialization here
	// ======================================

	// initialize aws config
	awscfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("error on loading aws config. err=%v", err)
	}

	// initialize postgres client
	connectionString := `0.0.0.0 user=local password=local dbname=image_poster port=5432 sslmode=disable`
	postgresClient := postgres.New(connectionString)

	// initialize repositories
	postRepo := postpkgrepo.New(postgresClient.GetDb())

	// initialize modules
	filestorageS3module := filestorageS3pkg.New(awscfg)
	postmodule := postpkg.New(postRepo, filestorageS3module)

	// NOTE: add server initialization on the following
	// ================================================
	posthandler := handlers.NewPost(postmodule)  // initialize handler
	httpserver := http.New(":4000", posthandler) // initialize http server

	// NOTE: add server starter on the following
	// =========================================
	postgresClient.Start() // start postgres connection
	httpserver.Start()     // start http server

	// "hang" the server and keep it
	// running until app got exit signal
	waitForExitSignal()

	// shutdown starts
	log.Println("shutting down server...")

	// NOTE: add shutdown handler on the following
	// ===========================================

	httpserver.Stop()     // stopping http server
	postgresClient.Stop() // closing postgress connection

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
