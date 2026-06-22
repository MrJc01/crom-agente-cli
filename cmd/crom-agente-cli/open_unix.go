//go:build !windows

package main

import (
	"os/exec"
	"syscall"
)

func startCommand(cmd *exec.Cmd) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	return cmd.Start()
}
