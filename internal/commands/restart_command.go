package commands

import (
	"errors"
	"fmt"
	"path/filepath"
	"runner/internal/services"
	"strconv"

	"github.com/urfave/cli/v3"
)

type RestartCommand struct {
	Cmd *cli.Command
}

func (Cmd RestartCommand) getPid(name string) (*int, error) {
	pids, err := services.GetPids(name)
	if err != nil {
		fmt.Println("Error getting PIDs (internal error)")
		return nil, errors.New("error getting PIDs")
	}
	if len(pids) == 1 {
		return pids[0], nil
	}

	if len(pids) == 0 {
		fmt.Printf("No process found with name \"%s\"\n", name)
		return nil, errors.New("no process found with that name")
	}

	if len(pids) > 1 {
		fmt.Printf("Multiple processes with name \"%s\" found\n", name)
		fmt.Printf("Use --pid for a specific process or --all to restart all matching \"%s\"\n", name)
		return nil, nil
	}
	return nil, nil
}

// TODO: all flag
func (Cmd RestartCommand) Run() error {
	if Cmd.Cmd.Args().Len() < 1 {
		fmt.Println("Please enter a valid process name")
	}
	allFlag := Cmd.Cmd.Bool("all")
	pidFlag := Cmd.Cmd.Bool("pid")
	var pid int
	name := Cmd.Cmd.Args().First()
	if pidFlag {
		p, err := strconv.Atoi(name)
		if err != nil {
			fmt.Println("Please enter a valid pid (numbers only)")
			return nil
		}
		pid = p
	} else {
		if !allFlag {
			if p, err := Cmd.getPid(name); err == nil {
				pid = *p
			} else {
				return nil
			}
		}
	}
	tmpDir := services.GetTempDirPath()
	activity, err := services.GetActivity(filepath.Join(tmpDir, strconv.Itoa(pid)+".json"))

	if err != nil {
		fmt.Printf("Process with id %s not found\n", name)
		return nil
	}

	err = services.StopActivity(pid)
	if err != nil {
		fmt.Println("Error stopping process (internal error)")
		return nil
	}

	_, err = services.StartProcessInBackground(append([]string{activity.Command}, activity.Arguments...)...)
	if err != nil {
		fmt.Println("Error starting process (internal error)")
		return nil
	}
	return nil
}
