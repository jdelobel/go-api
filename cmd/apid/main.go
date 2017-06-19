package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-api/cmd/apid/handlers"
	"github.com/go-api/internal/platform/db"
)

// init is called before main. We are using init to customize logging output.
func init() {
	log.SetFlags(logLstdFlags | log.Lmicroseconds | log.Lshortfile)
}

// main is the entry point for the application.
func main() {
	log.Println("main : Started")

	// Check the environment for a configured port value.
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres://images:images@127.0.0.1:54321/images?sslmode=disable"
	}

	// Register the Master Session for the database.
	log.Println("main : Started : Capturing Master DB...")

	masterDB, err := db.NewPSQL("postgres", dbHost)
	if err != nil {
		log.Fatalf("startup : Register DB : %v", err)
	}
	defer masterDB.PSQLClose()

	// Check the environment for a configured port value.
	host := os.Getenv("HOST")
	if host == "" {
		host = ":3000"
	}

	// Create a new server and set timeout values.
	server := http.Server{
		Addr:           host,
		Handler:        handlers.API(masterDB),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// We want to report the listener is closed.
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the listener.
	go func() {
		log.Printf("startup : Listening %s", host)
		log.Printf("shutdown : Listener closed : %v", server.ListenAndServe())
		wg.Done()
	}()

	// Listen for an interrupt signal from the OS.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt)

	// Wait for a signal to shutdown.
	<-osSignals

	// Create a context to attempt a graceful 5 second shutdown.
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Attempt the graceful shutdown by closing the listener and
	// completing all inflight requests.
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown : Graceful shutdown did not complete in %v : %v", timeout, err)

		// Looks like we timedout on the graceful shutdown. Kill it hard.
		if err := server.Close(); err != nil {
			log.Printf("shutdown : Error killing server : %v", err)
		}
	}

	// Wait for the listener to report it is closed.
	wg.Wait()
	log.Println("main : Completed")
}
