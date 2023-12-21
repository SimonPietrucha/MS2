package main

import (
	"context"
	"fmt"
	"time"

	"github.com/SimonPietrucha/MS2/Anwendung"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	app, err := Anwendung.New()
	if err != nil {
		// Behandle den Fehler
		fmt.Println("Failed to create app:", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = app.DB.Client().Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return
	}

	fmt.Println("Successfully connected and pinged MongoDB")

	err = app.Start(context.TODO())
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
