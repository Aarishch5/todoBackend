package main

import (
	"ToDo/database/migrations"
	"ToDo/server"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	err := migrations.OpenConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer migrations.DB.Close()

	myServer := server.SetupRoutes()

	fmt.Println("Server Listening at 8080")
	myServer.Start(":8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := myServer.Shutdown(ctx); err != nil {
		log.Fatalf("forcefully shutdown: %v", err)
	}

	log.Println("Server exited")
}
