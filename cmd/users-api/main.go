package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Divyansh031/goProject/internal/config"
)

func main() {
	//load config
	cfg  := config.MustLoad()

	//database setup

	//router setup
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to users-api!"))
	})


	// server setup
	server := http.Server{
		Addr: cfg.HTTPServer.Addr,
		Handler: router,
	}

	slog.Info("Server Started",slog.String("address", cfg.HTTPServer.Addr))


	done := make(chan os.Signal, 1) // to catch interrupt  signals

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func(){
		err := server.ListenAndServe()
		if err != nil{
			log.Fatal("Failed to start server")
			}

	}()
	

	<- done  // wait for interrupt signal
    
	slog.Info("Shutting down the server") // Logging
	
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil{
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
		
	}

	slog.Info("Server shutdown successfully")
	
}