package commands

import (
	"fmt"
	"runner/internal/services"

	"github.com/urfave/cli/v3"
)

type RunCommand struct {
	Cmd *cli.Command
}

func (Cmd RunCommand) Run() error {
	if Cmd.Cmd.Args().Len() < 1 {
		fmt.Println("Must provide a command!")
		return nil
	}
	activity, err := services.StartProcessInBackground(Cmd.Cmd.Args().Slice()...)
	if err != nil {
		fmt.Println("Error while executing command")
		return nil
	}

	fmt.Printf("Started %s (pid %d) in the background!\n", activity.Command, activity.Pid)

	return nil
}
