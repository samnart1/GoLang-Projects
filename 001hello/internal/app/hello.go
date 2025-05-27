package app

import (
	"fmt"
	"log"
	"time"
)

type Config struct {
	Name		string
	Verbose		bool
}

type App struct {
	config 		Config
	log			*log.Logger
}

func New(config Config) *App {
	return &App{
		config: config,
		log: log.Default(),
	}
} 

func (a *App) Run() error {
	if a.config.Verbose {
		a.log.Printf("Starting the application at %s", time.Now().Format(time.RFC3339))
		a.log.Printf("Application configurations: Name:%s, Verbose:%t", a.config.Name, a.config.Verbose)
	}

	message := a.generateGreeting()
	fmt.Println(message)

	if a.config.Verbose {
		a.log.Printf("Application finished successfully!")
	}

	return nil
}

func (a *App) generateGreeting() string {
	return fmt.Sprintf("Hello, %s", a.config.Name)
}