package service

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mhs003/harbrix/internal/paths"
)

func (s *State) Start(paths *paths.Paths) error {
	if s.Running {
		return errors.New("service already running")
	}

	cmd := exec.Command("sh", "-c", s.Config.Service.Command)
	if s.Config.Service.Workdir != "" {
		cmd.Dir = s.Config.Service.Workdir
	}

	logFile := filepath.Join(paths.ServiceLogs, s.Config.Name+".log")
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	cmd.Stdout = f
	cmd.Stderr = f

	if err := cmd.Start(); err != nil {
		return err
	}

	s.Running = true
	s.PID = cmd.Process.Pid
	s.Cmd = cmd
	s.StopReq = false

	go s.wait(paths)

	return nil
}

func (s *State) Stop() error {
	if !s.Running || s.Cmd == nil {
		return errors.New("service not running")
	}

	s.StopReq = true

	if err := s.Cmd.Process.Kill(); err != nil {
		return err
	}

	s.Cmd = nil
	s.Running = false
	s.PID = 0

	return nil
}

func (s *State) wait(paths *paths.Paths) {
	err := s.Cmd.Wait()

	s.Running = false
	s.PID = 0

	exitCode := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		}
	}
	s.ExitCode = exitCode

	if s.StopReq {
		return
	}

	switch s.Config.Service.Restart {
	case "always":
		s.Start(paths)
	case "on-failure":
		if exitCode != 0 {
			s.Start(paths)
		}
	}
}
