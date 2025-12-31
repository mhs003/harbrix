package cli

import (
	"fmt"
	"path/filepath"

	"github.com/mhs003/harbrix/internal/paths"
)

func CmdList(_ []string) {
	resp := send("list", "")
	services := resp.Data["services"].([]any)

	fmt.Printf("%s%-15s %-25s %-10s%s\n", bold, "SERVICE", "STATUS", "ENABLED", reset)
	fmt.Println("-----------------------------------------------------------")

	for _, raw := range services {
		s := raw.(map[string]any)
		name := s["name"].(string)
		running := s["running"].(bool)
		pid := s["pid"].(float64)
		enabled := s["enabled"].(bool)

		statusText := "stopped"
		if running {
			statusText = fmt.Sprintf("running (pid %.0f)", pid)
		}

		statusStartColor := red
		if running {
			statusStartColor = green
		}

		// Enabled string
		en := fmt.Sprintf("%sno%s", bold, reset)
		if enabled {
			en = fmt.Sprintf("%s%syes%s", bold, blue, reset)
		}

		fmt.Printf("%-15s %s%-25s%s %-10s\n", name, statusStartColor, statusText, reset, en)
	}
}

func CmdReloadDaemon(args []string) {
	send("reload-daemon", "")
	Success("Daemon reloaded.")
}

func CmdStart(args []string) {
	requireArg(args)
	send("start", args[0])
	Info("Service started.")
}

func CmdStop(args []string) {
	requireArg(args)
	send("stop", args[0])
	Info("Service stopped.")
}

func CmdRestart(args []string) {
	requireArg(args)
	send("stop", args[0])
	send("start", args[0])
	Info("Service restarted.")
}

func CmdLog(args []string) {
	if len(args) == 0 {
		Fatal("service name required")
	}

	follow := false
	name := args[0]

	if args[0] == "-f" {
		if len(args) < 2 {
			Fatal("service name required")
		}
		follow = true
		name = args[1]
	}

	showLog(name, follow)
}

func CmdEnable(args []string) {
	requireArg(args)
	send("enable", args[0])

	if len(args) > 1 && args[1] == "--now" {
		send("start", args[0])
		Info("Service enabled and started.")
	} else {
		Info("Service enabled.")
	}
}

func CmdDisable(args []string) {
	requireArg(args)
	send("disable", args[0])
	Info("Service disabled.")
}

func CmdIsEnabled(args []string) {
	requireArg(args)
	resp := send("is-enabled", args[0])
	if resp.Ok {
		Warning("Service is enabled")
	}
	Info("Service is disabled")
}

func CmdNew(args []string) {
	requireArg(args)
	createService(args[0])
}

func CmdEdit(args []string) {
	requireArg(args)
	editService(args[0])
}

func CmdDelete(args []string) {
	requireArg(args)
	resp := send("delete", args[0])
	if resp.Ok {
		//
	}
}

func CmdStatus(args []string) {
	requireArg(args)
	name := args[0]

	p, err := paths.New()
	if err != nil {
		Fatal(err.Error())
	}

	resp := send("list", "")
	services := resp.Data["services"].([]any)

	var found bool
	for _, raw := range services {
		s := raw.(map[string]any)
		if s["name"].(string) == name {
			found = true

			running := s["running"].(bool)
			pid := s["pid"].(float64)
			enabled := s["enabled"].(bool)

			file_path := filepath.Join(p.Services, name+".toml")

			fmt.Printf("%sDaemon information and status%s\n", bold, reset)
			fmt.Printf("  %sService:%s %s\n", yellow, reset, name)
			fmt.Printf("  %sDescription:%s %s\n", yellow, reset, s["description"].(string))
			fmt.Printf("  %sAuthor:%s %s\n", yellow, reset, s["author"].(string))

			status := "stopped"
			statusColor := red
			if running {
				status = fmt.Sprintf("running %s%s(pid %.0f)", reset, blue, pid)
				statusColor = green
			}
			if enabled {
				status += fmt.Sprintf("%s%s%s (enabled)", reset, bold, green)
			}

			fmt.Printf("  %sStatus:%s %s%s%s\n", yellow, reset, statusColor, status, reset)
			fmt.Printf("  %sService File:%s %s\n", yellow, reset, file_path)
			break
		}
	}

	if !found {
		red := "\033[31m"
		reset := "\033[0m"
		fmt.Printf("%sService %s not found.%s\n", red, name, reset)
	}
}

func ShowHelp() {

}
