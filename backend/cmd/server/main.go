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

	"agriedge/backend/internal/adapter/mqtt"
	"agriedge/backend/internal/adapter/postgres"
	"agriedge/backend/internal/config"
	deliveryHTTP "agriedge/backend/internal/delivery/http"
	"agriedge/backend/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.Println("[Main] Starting AgriEdge Backend Service...")

	cfg := config.LoadConfig()

	// 1. Database Connection (TimescaleDB / PostgreSQL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("[DB] Warning: Unable to create connection pool: %v", err)
	} else if err := pgPool.Ping(ctx); err != nil {
		log.Printf("[DB] Warning: Unable to connect to TimescaleDB at %s: %v", cfg.DatabaseURL, err)
		log.Println("[DB] Note: Ensure TimescaleDB container is running via `docker-compose up -d`.")
	} else {
		log.Println("[DB] Successfully connected to TimescaleDB!")
	}
	if pgPool != nil {
		defer pgPool.Close()
	}

	// 2. Repositories and Use Cases
	pgRepo := postgres.NewRepository(pgPool)
	ingestUC := usecase.NewIngestMetricUseCase(pgRepo)
	queryMetricUC := usecase.NewQueryMetricsUseCase(pgRepo)

	// 3. MQTT Subscriber for ESP32 Nodes
	mqttSub := mqtt.NewSubscriber(cfg.MQTTBrokerURL, cfg.MQTTTopic, ingestUC)
	if err := mqttSub.Start(); err != nil {
		log.Printf("[MQTT] Warning: Could not connect to broker at %s: %v", cfg.MQTTBrokerURL, err)
		log.Println("[MQTT] Note: Ensure Mosquitto broker is running via `docker-compose up -d`.")
	} else {
		log.Println("[MQTT] Successfully connected and subscribed to MQTT broker!")
	}
	defer mqttSub.Stop()

	// 4. HTTP Server for Web Dashboard & REST API
	httpHandler := deliveryHTTP.NewHandler(queryMetricUC, pgRepo)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:      httpHandler.Routes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("[HTTP] AgriEdge REST API listening on port %s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[HTTP] Listen error: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("[Main] Shutting down service gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("[HTTP] Server forced to shutdown: %v", err)
	}

	log.Println("[Main] AgriEdge Backend stopped.")
}
