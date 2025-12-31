package helpers

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type User struct {
	Name  string
	UID   int
	GID   int
	Home  string
	Shell string
}

type LoginUserOptions struct {
	ExceptUID     int
	MinUID        int
	AllowedShells []string
	RequireHome   bool
}

func DefaultLoginUserOptions() LoginUserOptions {
	return LoginUserOptions{
		ExceptUID: 0,
		MinUID:    1000,
		AllowedShells: []string{
			"/bin/bash",
			"/bin/zsh",
			"/bin/sh",
			"/usr/bin/bash",
			"/usr/bin/zsh",
			"/usr/bin/sh",
			"/usr/bin/fish",
		},
		RequireHome: true,
	}
}

func GetLoginUsers(opts LoginUserOptions) ([]User, error) {
	f, err := os.Open("/etc/passwd")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	allowedShell := make(map[string]struct{}, len(opts.AllowedShells))
	for _, s := range opts.AllowedShells {
		allowedShell[s] = struct{}{}
	}

	var users []User
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 7 {
			continue
		}

		uid, err := strconv.Atoi(parts[2])
		if err != nil || (uid < opts.MinUID && uid != opts.ExceptUID) {
			continue
		}

		shell := parts[6]
		if _, ok := allowedShell[shell]; !ok {
			continue
		}

		home := parts[5]
		if opts.RequireHome {
			if _, err := os.Stat(home); err != nil {
				continue
			}
		}

		gid, _ := strconv.Atoi(parts[3])

		users = append(users, User{
			Name:  parts[0],
			UID:   uid,
			GID:   gid,
			Home:  home,
			Shell: shell,
		})
	}

	return users, scanner.Err()
}
