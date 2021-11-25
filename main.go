package main

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 1.2.0
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

	"github.com/project-jano/user-service.go/api"
	"github.com/project-jano/user-service.go/app"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	appConfiguration, err := app.LoadAppConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	mongoDBOptions := options.Client().
		ApplyURI(appConfiguration.DatabaseConfiguration.Uri).
		SetMaxPoolSize(appConfiguration.DatabaseConfiguration.MaxPoolSize).
		SetMinPoolSize(appConfiguration.DatabaseConfiguration.MinPoolSize)

	mongoDBClient, mongodbErr := mongo.NewClient(mongoDBOptions)
	if mongodbErr != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = mongoDBClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(mongoDBClient, ctx)

	// Check the connection to MongoDB
	err = mongoDBClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	janoAPI := api.NewAPI(appConfiguration.APIConfiguration, mongoDBClient)

	srv := &http.Server{
		Handler:      janoAPI.Router,
		Addr:         appConfiguration.ServerAddr,
		ReadTimeout:  appConfiguration.ServerReadTimeout,
		WriteTimeout: appConfiguration.ServerWriteTimeout,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server starting...")
		if appConfiguration.APIConfiguration.AuthEnabled {
			log.Println("API access is authenticated using Basic Authentication")
		} else {
			log.Println("API access does not require authentication")
		}

		if len(appConfiguration.ServerPrivateKeyFilePath) > 0 && len(appConfiguration.ServerCertificateFilePath) > 0 {
			log.Println("TLS enabled...")
			if err := srv.ListenAndServeTLS(appConfiguration.ServerCertificateFilePath, appConfiguration.ServerPrivateKeyFilePath); err != http.ErrServerClosed {
				log.Fatalf("ListenAndServe(): %v", err)
			}
			return
		}
		log.Println("TLS disabled...")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}

	}()

	<-done
	log.Print("Server stopped.")
}
