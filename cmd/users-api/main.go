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
	"github.com/Divyansh031/goProject/internal/http/handlers/users"
	"github.com/Divyansh031/goProject/internal/storage/sqlite"
)

func main() {
	//load config
	cfg  := config.MustLoad()

	//database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}

	slog.Info("Storage initialized successfully", slog.String("env", cfg.Env))




	//router setup
	router := http.NewServeMux()

	router.HandleFunc("POST /api/users", user.New(storage))
	router.HandleFunc("GET /api/users/{id}", user.GetUserById(storage))
	router.HandleFunc("GET /api/users", user.GetUserList(storage))


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
	
    

//    **What this does:**
// Start the server in a **separate goroutine** (background thread).

// **Why goroutine?**
// Because `server.ListenAndServe()` BLOCKS (waits forever, handling requests). If we didn't use `go`, the program would be stuck here and never reach the next line!
// ```
// Main Goroutine                    Server Goroutine
//       ↓                                  ↓
//   Start server in background    → Actually runs server
//       ↓                              (handles requests)
//   Keep going...
//       ↓
//   Wait for shutdown signal


	<- done  // wait for interrupt signal



// 	**THIS IS THE KEY LINE!**

// The program **BLOCKS HERE** and waits. It won't continue until something is sent to the `done` channel.

// **When does something get sent?**
// When you press **Ctrl+C**! The OS sends a signal to the channel.
// ```
// Program running...
// User presses Ctrl+C
//    ↓
// Signal sent to 'done' channel
//    ↓
// <- done receives the signal
//    ↓
// Program continues to shutdown code
    
	slog.Info("Shutting down the server") // Logging
	
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
	
}