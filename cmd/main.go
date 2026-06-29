package main

import (
	"ToDo/database/migrations"
	middleware "ToDo/middlewares"
	"fmt"
	"log"
	"net/http"

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

	fmt.Println("Server listening on :8080")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
