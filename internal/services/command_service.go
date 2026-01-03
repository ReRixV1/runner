package services

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ExecCommand(commands ...string) {
	cmd := exec.Command(commands[0], commands[1:]...)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("Error executing command")
		fmt.Println(err)
		return
	}

	fmt.Println(string(stdout))
}
func ExecCommandInBackground(commands ...string) error {
	var outb, errb bytes.Buffer
	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
