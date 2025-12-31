package main

import (
	"fmt"
	"os"

	"github.com/mhs003/harbrix/internal/cli"
)

type Command struct {
	Name string
	Args string
	Desc string
	Run  func(args []string)
}

var commands = []Command{
	{
		Name: "list",
		Desc: "List all services",
		Run:  cli.CmdList,
	},
	{
		Name: "start",
		Args: "<service>",
		Desc: "Start a service",
		Run:  cli.CmdStart,
	},
	{
		Name: "stop",
		Args: "<service>",
		Desc: "Stop a service",
		Run:  cli.CmdStop,
	},
	{
		Name: "restart",
		Args: "<service>",
		Desc: "Restart a service",
		Run:  cli.CmdRestart,
	},
	{
		Name: "log",
		Args: "[-f] <service>",
		Desc: "Show service logs",
		Run:  cli.CmdLog,
	},
	{
		Name: "enable",
		Args: "<service> [--now]",
		Desc: "Enable a service",
		Run:  cli.CmdEnable,
	},
	{
		Name: "disable",
		Args: "<service>",
		Desc: "Disable a service",
		Run:  cli.CmdDisable,
	},
	{
		Name: "is-enabled",
		Args: "<service>",
		Desc: "Check if service is enabled",
		Run:  cli.CmdIsEnabled,
	},
	{
		Name: "new",
		Args: "<service>",
		Desc: "Create new service file",
		Run:  cli.CmdNew,
	},
	{
		Name: "edit",
		Args: "<service>",
		Desc: "Edit service file",
		Run:  cli.CmdEdit,
	},
	{
		Name: "delete",
		Args: "<service>",
		Desc: "Delete a service if it is not enabled and running",
		Run:  cli.CmdDelete,
	},
	{
		Name: "reload-daemon",
		Desc: "Reload daemon service files",
		Run:  cli.CmdReloadDaemon,
	},
	{
		Name: "status",
		Args: "<service>",
		Desc: "Show status of a service",
		Run:  cli.CmdStatus,
	},
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	name := os.Args[1]
	args := os.Args[2:]

	for _, c := range commands {
		if c.Name == name {
			c.Run(args)
			return
		}
	}

	cli.Fatal("unknown command: " + name)
}

func printHelp() {
	fmt.Println("harbrix - service manager")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  harbrix <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")

	for _, c := range commands {
		line := fmt.Sprintf("  %-15s", c.Name)
		if c.Args != "" {
			line += " " + c.Args
		}
		fmt.Printf("%-32s	%s\n", line, c.Desc)
	}
}
