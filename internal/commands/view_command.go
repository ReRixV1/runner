package commands

import (
	"fmt"
	"path/filepath"
	"runner/internal/config"
	"runner/internal/services"
	"strconv"

	"github.com/urfave/cli/v3"
)

type ViewCommand struct {
	Cmd     *cli.Command
	UseTail bool
}

func (Cmd ViewCommand) Run() error {
	tmpDir := services.GetTempDirPath()

	lines := config.Cfg.StartLines
	lines = max(lines, 0)

	var pid int
	if Cmd.Cmd.Args().Len() == 0 {
		fmt.Println("Please enter a process name or a PID using --pid")
		return nil
	}
	if Cmd.Cmd.Bool("pid") {
		p, err := strconv.Atoi(Cmd.Cmd.Args().First())
		pid = p
		if err != nil {
			fmt.Println("Please enter a valid pid!")
			return nil
		}
	}

	pids, _ := services.GetPids(Cmd.Cmd.Args().First())
	name := Cmd.Cmd.Args().First()
	if len(pids) > 1 {
		fmt.Printf("More than one process found with name \"%s\"\n", name)
		return nil
	}

	if !Cmd.Cmd.Bool("pid") && len(pids) == 0 {
		fmt.Printf("Process \"%s\" not found\n", name)
		return nil
	}

	if !Cmd.Cmd.Bool("pid") {
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
	if Cmd.UseTail {
		services.ReadLogFileTail(logFilePath, lines)
		return nil
	}

	services.ReadLogFile(logFilePath)
	return nil

}
