package commands

import (
	"fmt"
	"runner/internal/services"

	"github.com/urfave/cli/v3"
)

type ListCommand struct {
	Cmd *cli.Command
}

func (Cmd ListCommand) Run() error {
	err := services.ListActivites()
	if err != nil {
		fmt.Println("Error listing activities (internal error)")
		return nil
	}
	return nil
}
