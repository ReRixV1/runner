package commands

import (
	"fmt"
	"path/filepath"
	"runner/internal/services"
	"strconv"

	"github.com/urfave/cli/v3"
)

type ViewCommand struct {
	Cmd     *cli.Command
	UseTail bool
}

func (V ViewCommand) Run() error {
	tmpDir := services.GetTempDirPath()
	var pid int
	if V.Cmd.Args().Len() == 0 {
		fmt.Println("Please enter a process name or a PID using --pid")
		return nil
	}
	if V.Cmd.Bool("pid") {
		p, err := strconv.Atoi(V.Cmd.Args().First())
		pid = p
		if err != nil {
			fmt.Println("Please enter a valid pid!")
			return nil
		}
	}

	pids, _ := services.GetPids(V.Cmd.Args().First())
	name := V.Cmd.Args().First()
	if len(pids) > 1 {
		fmt.Printf("More than one process found with name \"%s\"\n", name)
		return nil
	}

	if !V.Cmd.Bool("pid") && len(pids) == 0 {
		fmt.Printf("Process \"%s\" not found\n", name)
		return nil
	}

	if !V.Cmd.Bool("pid") {
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
	if V.UseTail {
		services.ReadLogFileTail(logFilePath)
		return nil
	}

	services.ReadLogFile(logFilePath)
	return nil

}
