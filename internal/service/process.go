package service

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/mhs003/harbrix/internal/paths"
)

func (s *State) Start(paths *paths.Paths) error {
	if s.Running {
		return errors.New("service already running")
	}

	if err := s.Config.ValidateConfig(); err != nil {
		return err
	}

	cmd := exec.Command("sh", "-c", s.Config.Service.Command)
	if s.Config.Service.Workdir != "" {
		cmd.Dir = s.Config.Service.Workdir
	} else {
		cmd.Dir = paths.Root
	}

	if s.Config.Service.Log {
		logFile := filepath.Join(paths.ServiceLogs, s.Config.Name+".log")
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}

		cmd.Stdout = f
		cmd.Stderr = f

		syscall.Fchown(int(f.Fd()), s.UID, s.GID)
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid:         uint32(s.UID),
			Gid:         uint32(s.GID),
			Groups:      []uint32{uint32(s.GID)},
			NoSetGroups: false,
		},
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	s.Running = true
	s.PID = cmd.Process.Pid
	s.Cmd = cmd
	s.StopReq = false
	s.LastStartTime = time.Now()

	go s.wait(paths)

	return nil
}

func (s *State) Stop() error {
	if !s.Running || s.Cmd == nil {
		return errors.New("service not running")
	}

	s.StopReq = true

	pgid, err := syscall.Getpgid(s.Cmd.Process.Pid)
	if err != nil {
		return errors.New(err.Error())
	}

	if err := syscall.Kill(-pgid, syscall.SIGTERM); err != nil {
		return errors.New(err.Error())
	}
	if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
		return errors.New(err.Error())
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

	if exitCode != 0 {
		s.FailedCount++
	} else {
		s.FailedCount = 0
	}

	if !s.shouldRestart(exitCode) {
		return
	}

	delay, _ := time.ParseDuration(s.Config.Restart.Delay)
	if delay > 0 {
		time.Sleep(delay)
	}

	s.RestartCount++
	_ = s.Start(paths)
}

func (s *State) shouldRestart(exitCode int) bool {
	cfg := s.Config

	switch s.Config.Restart.Policy {
	case "never":
		return false
	case "on-failure":
		if exitCode == 0 {
			return false
		}
		// case "always":
	}

	maxFailed := cfg.Restart.MaxFailed
	if maxFailed != -1 && s.FailedCount > maxFailed-1 {
		s.FailedCount = 0
		return false
	}

	restartLimit := cfg.Restart.Limit
	if restartLimit > 0 && s.RestartCount >= restartLimit-1 {
		s.RestartCount = 0
		return false
	}

	return true
}
