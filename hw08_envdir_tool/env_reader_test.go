package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	vars := map[string]EnvValue{}
	vars["FOO"] = EnvValue{Value: "bar"}
	vars["LEN0"] = EnvValue{NeedRemove: true}
	vars["EMPTY"] = EnvValue{Value: ""}
	vars["BUZZ"] = EnvValue{Value: "world\n"}
	tempDir, err := os.MkdirTemp("", "*")
	require.NoError(t, err)
	os.WriteFile(path.Join(tempDir, "FOO"), []byte(vars["FOO"].Value), 0644)
	os.WriteFile(path.Join(tempDir, "LEN0"), []byte{}, 0644)
	os.WriteFile(path.Join(tempDir, "EMPTY"), []byte(" \t\n"), 0644)
	os.WriteFile(path.Join(tempDir, "BUZZ"), []byte("world\x00"), 0644)
	env, err := ReadDir(tempDir)
	require.NoError(t, err)
	require.Equal(t, vars["FOO"], env["FOO"])
	require.Equal(t, vars["LEN0"], env["LEN0"])
	require.Equal(t, vars["EMPTY"], env["EMPTY"])
	require.Equal(t, vars["BUZZ"], env["BUZZ"])
}
