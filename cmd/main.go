package main

import (
	"ToDo/database/migrations"
	middleware "ToDo/middlewares"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ToDo/handlers"

	"github.com/go-chi/chi/v5"
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

	router := chi.NewRouter()

	router.Get("/", handlers.HomeRoute)

	// Public Routes
	router.Post("/register", handlers.Register)
	router.Post("/login", handlers.Login)

	// Protected Routes
	router.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware)

		protected.Post("/todo", handlers.CreateTodo)
		protected.Get("/todos", handlers.GetTodos)
		protected.Get("/todo/{id}", handlers.GetTodo)
		protected.Put("/todo/{id}", handlers.UpdateTodo)
		protected.Delete("/todo/{id}", handlers.DeleteTodo)
		protected.Delete("/user", handlers.DeleteUser)
	})

	myServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Running the server in the background so it doesn't block.
	go func() {
		fmt.Println("Server Listening at 8080")
		if err := myServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := myServer.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
