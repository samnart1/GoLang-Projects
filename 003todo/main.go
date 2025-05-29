package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/samnart1/GoLang-Projects/003todo/cmd"
	"github.com/samnart1/GoLang-Projects/003todo/internal/config"
	"github.com/samnart1/GoLang-Projects/003todo/pkg/version"
)

func main() {
	cfg := config.New()

	if len(os.Args) < 2 {
		cmd.ShowHelp()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "add", "a":
		if len(args) == 0 {
			fmt.Println("Error: Please projvide a task description")
			os.Exit(1)
		} 
		cmd.AddTask(cfg, args)

	case "list", "ls", "l":
		cmd.ListTasks(cfg, args)

	case "done", "complete", "d":
		if len(args) == 0 {
			fmt.Println("Error: Please provide a task ID")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: Invlid task ID")
			os.Exit(1)
		}
		cmd.CompleteTask(cfg, id)

	case "remove", "rm", "r":
		if len(args) == 0 {
			fmt.Println("Error: Please provide a task ID")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("error: invalid task ID")
			os.Exit(1)
		}
		cmd.RemoveTask(cfg, id)

	case "edit", "e":
		if len(args) < 2 {
			fmt.Println("Error: please provide task ID and new description")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: Invlaid task ID")
			os.Exit(1)
		}
		cmd.EditTask(cfg, id, args[1:])

	case "search", "s":
		if len(args) == 0 {
			fmt.Println("Error: Please provide a search term")
			os.Exit(1)
		}
		cmd.SearchTasks(cfg, args[0])

	case "stats":
		cmd.ShowStat(cfg)

	case "version", "v":
		fmt.Printf("todo-cli version %s\n", version.Version)

	case "help", "h":
		cmd.ShowHelp()

	default:
		fmt.Printf("Umknown command: %s\n", command)
		cmd.ShowHelp()
		os.Exit(1)
	}
}