package main

import (
	"context"
	"log"

	"Subscription_Service/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Printf("Failed to create a new app: %s", err.Error())
		return
	}

	err = app.Start(context.Background())
	if err != nil {
		log.Printf("Failed to start the app: %s", err.Error())
	}
}
