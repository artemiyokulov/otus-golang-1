package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	for key, value := range env {
		os.Unsetenv(key)
		if !value.NeedRemove {
			os.Setenv(key, value.Value)
		}
	}

	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Run()
	return command.ProcessState.ExitCode()
}
