package main

import (
	"fmt"
	"path/filepath"
	"runner/internal/services"
	"strconv"

	"github.com/urfave/cli/v3"
)

func view(c *cli.Command, useTail bool) error {
	tmpDir := services.GetTempDirPath()
	var pid int
	if c.Args().Len() == 0 {
		fmt.Println("Please enter a process name or a PID using --pid")
		return nil
	}
	if c.Bool("pid") {
		p, err := strconv.Atoi(c.Args().First())
		pid = p
		if err != nil {
			fmt.Println("Please enter a valid pid!")
			return nil
		}
	}

	pids, _ := services.GetPids(c.Args().First())

	if len(pids) > 1 {
		fmt.Printf("More than one process found with name \"%s\"\n", c.Args().First())
		return nil
	}

	if !c.Bool("pid") && len(pids) == 0 {
		fmt.Printf("Process \"%s\" not found\n", c.Args().First())
		return nil
	}

	if !c.Bool("pid") {
		pid = *pids[0]
	}

	path := filepath.Join(tmpDir, strconv.Itoa(pid)+".json")
	a, err := services.GetActivity(path)
	if err != nil {
		fmt.Println("pid not found!")
		return nil
	}

	logFile := a.LogFile
	logFilePath := filepath.Join(tmpDir, logFile)
	if useTail {
		services.ReadLogFileTail(logFilePath)
		return nil
	}

	services.ReadLogFile(logFilePath)
	return nil

}
