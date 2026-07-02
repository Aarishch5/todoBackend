package server

import (
	"context"
	"net/http"
	"time"

	"ToDo/handlers"
	middleware "ToDo/middlewares"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

// SetupRoutes provides all the routes that can be used
func SetupRoutes() *Server {
	router := chi.NewRouter()

	router.Route("/v1", func(v1 chi.Router) {
		v1.Get("/home", handlers.Home)

		v1.Post("/register", handlers.Register)
		v1.Post("/login", handlers.Login)

		v1.Group(func(protected chi.Router) {
			protected.Use(middleware.AuthMiddleware)

			protected.Post("/todo", handlers.CreateTodo)
			protected.Get("/todos", handlers.GetAllTodos)
			protected.Get("/todo/{id}", handlers.GetTodoById)
			protected.Put("/todo/{id}", handlers.UpdateTodoById)
			protected.Delete("/todo/{id}", handlers.DeleteTodoById)
			protected.Delete("/user", handlers.DeleteUser)
			protected.Post("/logout", handlers.Logout)
		})
	})

	return &Server{
		Router: router,
	}
}

func (s *Server) Start(addr string) {
	s.server = &http.Server{
		Addr:              addr,
		Handler:           s.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
