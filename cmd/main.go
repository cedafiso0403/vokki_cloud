package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"vokki_cloud/config"
	"vokki_cloud/docs" // This is required for Swagger to find your docs
	"vokki_cloud/internal/database"
	"vokki_cloud/internal/router"
	"vokki_cloud/internal/shared"
	"vokki_cloud/internal/constants"

	_ "github.com/lib/pq"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	// Create a WaitGroup to keep track of running goroutines
	var wg sync.WaitGroup

	wg.Add(1)

	server, logFile, err := config.LoadConfig()

	// Load the server configuration
	if err != nil {
		log.Fatal("Error loading config: ", err)
		return
	}

	defer logFile.Close() // Ensure the log file is closed when the program exits

	go startServer(ctx, &wg, server)

	// Listen for termination signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-signalCh

	// Start the graceful shutdown process
	log.Println("Gracefully shutting down HTTP server...")

	// Cancel the context to signal the HTTP server to stop
	cancel()

	// Wait for the HTTP server to finish
	wg.Wait()

	log.Println("Shutdown complete.")
	log.Println("Server stopped")
}

func startServer(ctx context.Context, wg *sync.WaitGroup, server *http.Server) {
	defer wg.Done()

	// Swagger documentation setup
	docs.SwaggerInfo.Title = "Vokki Cloud API"
	docs.SwaggerInfo.Description = "This is the API for Vokki mobile app."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = vokki_constants.BaseHost    
    docs.SwaggerInfo.BasePath = vokki_constants.BasePath 
    docs.SwaggerInfo.Schemes = []string{vokki_constants.BaseScheme} 

	// Initialize token manager
	shared.InitializeTokenManager()

	// Connect to the database and initialize prepared statements
	database.Connect()
	defer database.Close() // Ensure resources are cleaned up when main exits

	// Setup the router
	r := router.SetupRouter()
	server.Handler = r

	// Start the server in the current goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	select {
	case <-ctx.Done():
		// Shutdown the server gracefully
		log.Println("Shutting down HTTP server gracefully...")
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("HTTP server shutdown error: %s\n", err)
		}
		log.Println("HTTP server stopped.")
	}

}
