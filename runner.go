package main

import (
	"log"
	"os"
	"os/exec"
)

var runningCommands []*exec.Cmd

func kill() {
	for _, cmd := range runningCommands {
		err := cmd.Process.Kill()
		if err != nil {
			println(err)
		}
	}
}

func runCommand(commandText string) error {
	cmd := exec.Command("sh", "-c", commandText)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Println("$ " + commandText)
	err := cmd.Start()
	if err != nil {
		return err
	}
	runningCommands = append(runningCommands, cmd)
	return nil
}
