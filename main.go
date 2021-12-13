package main

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com
 */

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/project-jano/user-service.go/helpers"

	"github.com/project-jano/user-service.go/api"
	"github.com/project-jano/user-service.go/app"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	config, err := app.LoadAppConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	mongoDBClient := initMongoDB(config)

	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(mongoDBClient, context.Background())

	fingerprint := helpers.HashedFingerprint()

	janoAPI := api.NewAPI(config.API, mongoDBClient, fingerprint)

	srv := &http.Server{
		Handler:      janoAPI.Router,
		Addr:         config.Server.Addr,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.ServerWriteTimeout,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		startServer(config, srv)
	}()

	<-done
	log.Print("Server stopped.")
}

func startServer(appConfiguration *app.Configuration, srv *http.Server) {
	log.Println("Server starting...")
	if appConfiguration.API.AuthEnabled {
		log.Println("API access is authenticated using Basic Authentication")
	} else {
		log.Println("API access does not require authentication")
	}

	if len(appConfiguration.Server.PrivateKeyFilePath) > 0 && len(appConfiguration.Server.CertificateFilePath) > 0 {
		log.Println("TLS enabled...")
		if err := srv.ListenAndServeTLS(appConfiguration.Server.CertificateFilePath, appConfiguration.Server.PrivateKeyFilePath); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
		return
	}
	log.Println("TLS disabled...")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
}

func initMongoDB(appConfiguration *app.Configuration) *mongo.Client {
	mongoDBOptions := options.Client().
		ApplyURI(appConfiguration.Database.Uri).
		SetMaxPoolSize(appConfiguration.Database.MaxPoolSize).
		SetMinPoolSize(appConfiguration.Database.MinPoolSize)

	mongoDBClient, mongodbErr := mongo.NewClient(mongoDBOptions)
	if mongodbErr != nil {
		log.Fatal(mongodbErr)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := mongoDBClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection to MongoDB
	err = mongoDBClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return mongoDBClient
}
