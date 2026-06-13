// @title Todo App API
// @version 1.0.0
// @description A full-featured Todo application with JWT auth, Redis caching, and Prometheus metrics.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token.

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kaanchinar/todo-app/config"
	"github.com/kaanchinar/todo-app/handlers"
	apimw "github.com/kaanchinar/todo-app/middleware"
	"github.com/kaanchinar/todo-app/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"


	"github.com/kaanchinar/todo-app/docs"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s, err := store.NewCompositeStore(ctx, cfg.DatabaseURL, cfg.RedisAddr, cfg.RedisPass)
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}
	defer s.Close()

	authHandler := handlers.NewAuthHandler(s, cfg)
	todoHandler := handlers.NewTodoHandler(s)

	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(apimw.Metrics)

	// Public routes
	r.Get("/", handlers.Root(cfg))
	r.Get("/health", handlers.Health())
	r.Handle("/metrics", promhttp.Handler())

	// Scalar API docs
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(scalarHTML))
	})
	r.Get("/docs/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(docs.GetSwaggerJSON())
	})

	// Auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// Protected todo routes
	r.Route("/todos", func(r chi.Router) {
		r.Use(apimw.Auth(cfg))
		r.Get("/", todoHandler.List)
		r.Post("/", todoHandler.Create)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", todoHandler.Get)
			r.Put("/", todoHandler.Update)
			r.Delete("/", todoHandler.Delete)
		})
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		srv.Shutdown(shutdownCtx)
	}()

	log.Printf("🚀 %s v%s listening on %s", cfg.AppName, cfg.AppVersion, addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

const scalarHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Todo App API - Docs</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <link rel="icon" type="image/svg+xml" href="https://cdn.jsdelivr.net/npm/@scalar/api-reference/favicon.svg" />
    <style>
        body { margin: 0; }
    </style>
</head>
<body>
    <script id="api-reference"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    <script>
        document.getElementById('api-reference').configure({
            spec: {
                url: '/docs/openapi.json'
            }
        });
    </script>
</body>
</html>`
