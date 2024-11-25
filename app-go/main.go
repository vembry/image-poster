package main

import (
	"app-go/internal/app"
	"app-go/internal/clients/postgres"
	filestorageS3pkg "app-go/internal/modules/file_storage/s3"
	postpkgrepo "app-go/internal/modules/post/repositories/postgres"
	postpkg "app-go/internal/modules/post/services"
	"app-go/internal/servers/http"
	"app-go/internal/servers/http/handlers"
	"app-go/internal/workers/sqs"
	"context"
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

var (

	//go:embed configs
	embedFS embed.FS
)

func main() {
	log.Println("initiating pre-requisites...")

	// NOTE: add base initialization here
	// ======================================
	appConfig := app.NewConfig(embedFS)

	// initialize aws config
	awscfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("error on loading aws config. err=%v", err)
	}

	// initialize worker handler
	imageTransformWorkerSqs := sqs.NewImageTransformer(appConfig.AWS.Sqs.ImageTransformQueueUrl)

	// initialize worker engine with sqs
	sqsworker := sqs.New(awscfg)
	sqsworker.RegisterHandlers(
		imageTransformWorkerSqs,
		// NOTE: register more worker here
		// ...
	)

	// initialize postgres client
	postgresClient := postgres.New(appConfig.Postgres.ConnectionString)

	// initialize repositories
	postRepo := postpkgrepo.NewPost(postgresClient)
	postStructureRepo := postpkgrepo.NewPostStructure(postgresClient)

	// initialize modules
	filestorageS3module := filestorageS3pkg.New(awscfg)
	postmodule := postpkg.New(postRepo, postStructureRepo, filestorageS3module, imageTransformWorkerSqs)

	// NOTE: add server initialization on the following
	// ================================================
	posthandler := handlers.NewPost(postmodule)  // initialize handler
	httpserver := http.New(":4000", posthandler) // initialize http server

	// NOTE: add server starter on the following
	// =========================================
	log.Println("starting server/worker...")

	postgresClient.Start() // starts postgres connection
	httpserver.Start()     // starts http server
	sqsworker.Start()      // starts sqs worker

	// "hang" the server and keep it
	// running until app got exit signal
	waitForExitSignal()

	// shutdown starts
	log.Println("shutting down server...")

	// NOTE: add shutdown handler on the following
	// ===========================================

	postgresClient.Stop() // closes postgres connection
	httpserver.Stop()     // stops http server
	sqsworker.Stop()      // stops sqs worker

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
