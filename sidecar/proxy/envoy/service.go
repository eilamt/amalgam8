package envoy

import (
	"fmt"
	"os/exec"
)

type Service interface {
	Start() error
	Reload() error
	Stop() error
}

type service struct {
	cmd          *exec.Cmd
	config       string
	restartEpoch int
}

func NewService(config string) Service {
	return &service{
		config: config,
	}
}

func (s *service) Start() error {
	if s.cmd == nil {
		return nil
	}

	s.cmd = exec.Command("envoy", "-c", s.config, "--restart-epoch", fmt.Sprint(s.restartEpoch))
	return s.cmd.Start() // Wait() to prevent zombies?!
}

func (s *service) Reload() error {
	if s.cmd != nil {
		s.cmd.Process.Release()
	}

	s.restartEpoch++
	s.cmd = exec.Command("envoy", "-c", s.config, "--restart-epoch", fmt.Sprint(s.restartEpoch))
	return s.cmd.Start() // Wait() to prevent zombies?!
}

func (s *service) Stop() error {
	return s.cmd.Process.Kill()
}
