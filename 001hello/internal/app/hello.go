package app

import (
	"fmt"
	"log"
	"time"
)

type Config struct {
	Name string
	Verbose bool
}

type App struct {
	config Config
	logger *log.Logger
}

func New(config Config) *App {
	return &App{
		config: config,
		logger: log.Default(),
	}
}

func (a *App) Run() error {
	if a.config.Verbose {
		a.logger.Printf("Starting Hello CLI application at %s", time.Now().Format(time.RFC3339))
		a.logger.Printf("Configuration: Name=%s, Verbose=%t", a.config.Name, a.config.Verbose)
	}

	message := a.generateGreeting()
	fmt.Println(message)

	if a.config.Verbose{
		a.logger.Printf("Application completed successfully!")
	}

	return nil
}

func (a *App) generateGreeting() string {
	return fmt.Sprintf("Hello, %s!", a.config.Name)
}