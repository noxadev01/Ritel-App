package main

import (
	"context"
	"embed"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	customlogger "ritel-app/internal/logger"
	"ritel-app/internal/config"
	"ritel-app/internal/container"
	"ritel-app/internal/database"
	httpserver "ritel-app/internal/http"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Re-enable logging for initialization
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("========================================")
	log.Println("🚀 RITEL-APP STARTING")
	log.Println("========================================")

	// STEP 1: Initialize database FIRST
	log.Println("[INIT] Initializing database...")
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("[INIT] Database initialized successfully")

	// STEP 2: Create SERVICE CONTAINER (shared by both Wails & HTTP)
	log.Println("[INIT] Initializing service container...")
	services := container.NewServiceContainer()
	log.Println("[INIT] Service container initialized")

	// STEP 3: Create Wails app and inject services
	app := NewApp()
	app.SetServices(services)

	// STEP 4: Check if web server is enabled
	serverConfig := config.GetServerConfig()

	var httpServer *httpserver.Server
	if serverConfig.Enabled {
		log.Println("========================================")
		log.Println("🌐 WEB SERVER ENABLED")
		log.Println("========================================")

		// Create HTTP server
		httpServer = httpserver.NewServer(services, serverConfig)

		// Start HTTP server in background
		if err := httpServer.Start(); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}

		// Setup graceful shutdown for HTTP server
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := httpServer.Shutdown(ctx); err != nil {
				log.Printf("HTTP server shutdown error: %v", err)
			}
		}()

		// Handle OS signals for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigChan
			log.Println("\n[SHUTDOWN] Received shutdown signal")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if httpServer != nil {
				httpServer.Shutdown(ctx)
			}
			services.Shutdown()
			os.Exit(0)
		}()
	}

	// STEP 5: Run Wails application (BLOCKING)
	log.Println("========================================")
	log.Println("🖥️  STARTING DESKTOP APP (WAILS)")
	log.Println("========================================")

	// Disable logs for Wails runtime
	log.SetOutput(io.Discard)

	err := wails.Run(&options.App{
		Title:  "Ritel-App",
		Width:  1200,
		Height: 700,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown: func(ctx context.Context) {
			app.shutdown(ctx)
			services.Shutdown()
		},
		Bind: []interface{}{
			app,
		},
		Logger: customlogger.NewSilentLogger(),
	})

	if err != nil {
		log.SetOutput(os.Stdout)
		log.Printf("Error: %v", err)
	}

	log.SetOutput(os.Stdout)
	log.Println("========================================")
	log.Println("👋 APPLICATION SHUTDOWN COMPLETE")
	log.Println("========================================")
}
