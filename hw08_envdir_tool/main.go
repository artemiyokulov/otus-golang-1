package main

import (
	"log"
	"os"
)

func main() {
	env, err := ReadDir(os.Args[1])
	if err == nil {
		returnCode := RunCmd(os.Args[2:], env)
		os.Exit(returnCode)
	} else {
		log.Fatal(err)
	}
}
