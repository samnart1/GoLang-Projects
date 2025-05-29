package cmd

import (
	"fmt"

	"github.com/samnart1/GoLang-Projects/003todo/pkg/version"
)

func ShowHelp() {
	help := `Todo CLI - A simple command-line todo list manager

	USAGE:
		todo <command> [arguments]

	COMMANDS:
		add, a <description>		Add a new task
		list, ls, l [options]		List tasks
		done, complete, d <id>		Mark task as completed
		remove, rm, r <id>			Remove a task
		edit, e <id> <description>	Edit a task
		search, s <term>			Search tasks
		stats						Show statistics
		version, v					Show version
		help, h						Show this help

	LIST OPTIONS:
		--all						Show all tasks (default)
		--pending					Show only pending tasks
		--completed					Show only completed tasks
		--priority <level>			Filter by priority (low, medium, high)
		--table						Display in table format
	
	EXAMPLES:
		todo add "Buy groceries"
		todo add "Fix buy #123" --priority high
		todo list --pending
		todo list --table
		todo done 1
		todo remove 2
		todo edit 1 "Buy groceries and cook dinner"
		todo search "groceries"

	VERSION:
		%s
	`
	fmt.Printf(help, version.GetVersion())
}