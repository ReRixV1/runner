package commands

import (
	"fmt"
	"runner/internal/services"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"
)

type StopCommand struct {
	Cmd *cli.Command
}

func (S StopCommand) Run() error {
	services.DeleteStoppedActivites()

	if S.Cmd.Args().Len() < 1 {
		fmt.Println("Please enter a valid pid!")
		return nil
	}

	if S.Cmd.Bool("pid") == true {
		pid, err := strconv.Atoi(S.Cmd.Args().First())
		if err != nil {
			fmt.Println("Please enter a valid pid!")
			return nil
		}
		err = services.StopActivity(pid)

		if err != nil {
			fmt.Println("pid not found!")
			return nil
		}

		fmt.Println("Stopped process " + strconv.Itoa(pid))
		return nil
	}

	name := strings.ToLower(S.Cmd.Args().First())
	err := services.StopActivityWithName(name, S.Cmd.Bool("all"))

	if err != nil {
		if err.Error() == "not found" {
			return nil
		}
		fmt.Println("Error while trying to stop process!")
		return nil
	}

	fmt.Println("Stopped process: " + name)
	return nil
}
