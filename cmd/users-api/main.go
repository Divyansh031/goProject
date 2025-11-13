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
	user "github.com/Divyansh031/goProject/internal/http/handlers/users"
)

func main() {
	//load config
	cfg  := config.MustLoad()

	//database setup

	//router setup
	router := http.NewServeMux()

	router.HandleFunc("POST /api/users", user.New())


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

	err := server.Shutdown(ctx)
	if err != nil{
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
		
	}

	slog.Info("Server shutdown successfully")
	
}