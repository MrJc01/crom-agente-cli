//go:build windows

package main

import (
	"os/exec"
)

func startCommand(cmd *exec.Cmd) error {
	return cmd.Start()
}
