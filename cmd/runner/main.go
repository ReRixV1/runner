package main

import (
	"context"
	"log"
	"os"
	"runner/internal/commands"
	"runner/internal/services"

	"github.com/urfave/cli/v3"
)

func main() {
	if err := services.EnsureTempDirectory(); err != nil {
		return
	}

	cmd := &cli.Command{
		Name:                  "runner",
		Usage:                 "Manage commands that run in the background",
		Version:               "v0.1",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			&cli.Command{
				Name:            "run",
				Aliases:         []string{"r", "start"},
				Usage:           "Run command in the background",
				SkipFlagParsing: true,
				Action: func(ctx context.Context, c *cli.Command) error {
					return commands.RunCommand{Cmd: c}.Run()
				},
			},
			&cli.Command{
				Name:    "list",
				Aliases: []string{"l", "ls"},
				Usage:   "Lists all running backgrund activities (commands)",
				Action: func(ctx context.Context, c *cli.Command) error {
					return commands.ListCommand{Cmd: c}.Run()
				},
			},
			&cli.Command{
				Name:    "view",
				Aliases: []string{"v", "show"},
				Usage:   "View live output of process (experimental)",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "pid",
						Usage:   "view process using its pid",
						Aliases: []string{"p"},
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return commands.ViewCommand{Cmd: c, UseTail: true}.Run()
				},
			},
			&cli.Command{
				Name:  "log",
				Usage: "Show output of process",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "pid",
						Usage:   "show process output using its pid",
						Aliases: []string{"p"},
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return commands.ViewCommand{Cmd: c, UseTail: false}.Run()
				},
			},

			&cli.Command{
				Name:    "stop",
				Aliases: []string{"end", "kill", "s"},
				Usage:   "Stops an activity",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "pid",
						Usage:   "stop process using its pid",
						Aliases: []string{"p"},
					},
					&cli.BoolFlag{
						Name:    "all",
						Usage:   "stop all process matching the specified name",
						Aliases: []string{"a"},
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return commands.StopCommand{Cmd: c}.Run()
				},
			},
			&cli.Command{
				Name:  "restart",
				Usage: "restarts a process",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "pid",
						Usage:   "restart process using its pid",
						Aliases: []string{"p"},
					},
					&cli.BoolFlag{
						Name:    "all",
						Usage:   "restart all process matching name",
						Aliases: []string{"a"},
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return commands.RestartCommand{Cmd: c}.Run()
				},
			},
			//&cli.Command{
			//	Name:  "temp",
			//	Usage: "prints temp directory path (for development :D)",
			//	Action: func(ctx context.Context, c *cli.Command) error {
			//		fmt.Println(services.GetTempDirPath())
			//		return nil
			//	},
			//},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
