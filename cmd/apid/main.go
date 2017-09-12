package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"fmt"

	"github.com/jdelobel/go-api/cmd/apid/handlers"
	"github.com/jdelobel/go-api/config"
	"github.com/jdelobel/go-api/internal/platform/db"
	"github.com/jdelobel/go-api/internal/platform/rabbitmq"
	"github.com/jdelobel/go-api/logger"
)

// init is called before main. We are using init to customize logging output.
func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

// main is the entry point for the application.
func main() {
	log.Println("main : Started")

	var c = config.Config{}
	err := c.Load("config/config.json")
	if err != nil {
		log.Fatalf("main: Failed to init config: %v", err)
	}
	loggerConf := logger.Conf{Host: c.Logger.Host,
		Port:    c.Logger.Port,
		Level:   c.Logger.Level,
		App:     c.AppName,
		Version: c.AppVersion}
	if err = logger.Init(loggerConf); err != nil {
		log.Fatalf("main: Failed to init logger: %v", err)
	}
	dbHost := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		c.Database.Client,
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name)

	// Register the Master Session for the database.
	log.Println("main : Started : Capturing Master DB...")

	masterDB, err := db.NewPSQL(dbHost)
	if err != nil {
		log.Fatalf("startup : Register DB : %v", err)
	}

	rbmqHost := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		c.RabbitMQ.User,
		c.RabbitMQ.Password,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
		c.RabbitMQ.Name)
	defaultQueue := "go-api-messages"
	fmt.Println(rbmqHost)
	rbmq, err := rabbitmq.NewRabbitMQ(rbmqHost, &defaultQueue)

	if err != nil {
		log.Fatalf("startup : Register RabitMQ : %v", err)
	}
	host := fmt.Sprintf("%s:%s", c.AppHost, c.AppPort)
	// Create a new server and set timeout values.
	server := http.Server{
		Addr:           host,
		Handler:        handlers.API(masterDB, logger.Log, c, rbmq),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// We want to report the listener is closed.
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the listener.
	go func() {
		logger.Log.Infof("startup : Listening %s", host)
		logger.Log.Infof("shutdown : Listener closed : %v", server.ListenAndServe())
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
		logger.Log.Infof("shutdown : Graceful shutdown did not complete in %v : %v", timeout, err)

		// Looks like we timedout on the graceful shutdown. Kill it hard.
		if err := server.Close(); err != nil {
			logger.Log.Infof("shutdown : Error killing server : %v", err)
		}
	}
	if err := masterDB.PSQLClose(); err != nil {
		logger.Log.Errorf("main : Database instance not closed : %v", err)
	}

	// Wait for the listener to report it is closed.
	wg.Wait()
	logger.Log.Info("main : Completed")
}
