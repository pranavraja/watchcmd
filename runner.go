package main

import (
	"log"
	"os"
	"os/exec"
)

func runCommand(commandText string) error {
	cmd := exec.Command("sh", "-c", commandText)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Println("$ " + commandText)
	return cmd.Run()
}
