package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed web/public/*
var fileSys embed.FS

func main() {

	// The files in web/public/ shall be served at the base URL,
	// so we need to creat a sub-filesystem:

	publicFS, err := fs.Sub(fileSys, "web/public")
	if err != nil {
		log.Fatalf("fs.Sub: %v\n", err)
	}

	// Setup your multiplexer and handlers
	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServerFS(publicFS))

	// Configure timeouts! This is crucial.
	// The default http.Server doesn't do this.
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("Server starting on http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown logic
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop // Block until we receive a signal

	log.Println("Shutdown signal received, initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Shutdown timeout
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown due to error: %v", err)
	}
	log.Println("Server exited gracefully")
}
