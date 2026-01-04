package services

import (
	"fmt"
	"os"
	"os/exec"
)

func ReadLogFile(path string) error {
	cmd := exec.Command("tail", "-f", path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("Error reading log file (internal error)")
		return err
	}

	return nil
}
