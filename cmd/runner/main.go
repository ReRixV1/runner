package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runner/internal/services"

	"github.com/urfave/cli/v3"
)

func main() {
	if err := services.EnsureTempDirectory(); err != nil {
		return
	}

	cmd := &cli.Command{
		Name:    "runner",
		Usage:   "Manage commands that run in the background",
		Version: "v0.1",
		Commands: []*cli.Command{
			&cli.Command{
				Name:            "run",
				Aliases:         []string{"r"},
				Usage:           "Run command in the background",
				SkipFlagParsing: true,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						fmt.Println("Must provide command!")
						return nil
					}
					err := services.ExecCommandInBackground(cmd.Args().Slice()...)
					if err != nil {
						fmt.Println("Error while executing command")
						return nil
					}
					return nil
				},
			},
			&cli.Command{
				Name:    "list",
				Aliases: []string{"l", "ls"},
				Usage:   "Lists all running backgrund activites (commands)",
				Action: func(ctx context.Context, c *cli.Command) error {
					err := services.ListActivites()
					if err != nil {
						fmt.Println("Error listing activities (internal error)")
						return nil
					}
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
