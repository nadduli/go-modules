package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nadduli/go-modules/internal/config"
	"github.com/nadduli/go-modules/internal/db"
	"github.com/nadduli/go-modules/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}

	router := server.NewRouter(database)

	srv := &http.Server{
		Addr:        cfg.ServerPort,
		Handler:     router,
		ReadTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Printf("Error getting underlying DB connection: %v", err)
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing DB: %v", err)
		}
	}

	log.Println("Server exiting")
}
