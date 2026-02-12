package main

//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/ariashabry/boilerplate-go/helpers/log"
	"github.com/ariashabry/boilerplate-go/internal/domains/product/model"
)

func main() {
	// Initialize log
	l := log.NewLog("Log")
	defer l.Close()
	l.Info("Set Up Log Started")

	// Initialize migrations
	l.Info("Running database migrations")
	migrationService := InitializeMigrations(l)

	// Register models from all domains
	model.RegisterProductModels(migrationService)

	if err := migrationService.RunMigrations(); err != nil {
		l.WithError(err).Fatal("Failed to run database migrations")
	}

	// Wire everything up
	http := InitializeService(l)

	// Start HTTP server
	if err := http.SetupAndServe(); err != nil {
		l.Fatal(err)
	}
}
