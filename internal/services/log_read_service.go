package services

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func ReadLogFileTail(path string, lines int) error {
	cmd := exec.Command("tail", "-n", strconv.Itoa(lines), "-f", path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("Error reading log file (internal error)")
		return err
	}
	defer cmd.Process.Release()
	cmd.Wait()
	return nil
}

func ReadLogFile(path string) error {
	d, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error raeding log file (internal error)")
		return err
	}

	println(string(d))

	return nil
}
