package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := Environment{
		"FOO": EnvValue{Value: "bar"},
	}
	var exitCode int = RunCmd([]string{"/bin/sh", "-c", "echo $FOO"}, env)
	require.Equal(t, 0, exitCode)

	exitCode = RunCmd([]string{"false"}, Environment{})
	require.Equal(t, 1, exitCode)
}
