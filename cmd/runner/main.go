package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runner/internal/services"
	"strconv"
	"strings"

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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() < 1 {
						fmt.Println("Must provide a command!")
						return nil
					}
					activity, err := services.StartProcessInBackground(cmd.Args().Slice()...)
					if err != nil {
						fmt.Println("Error while executing command")
						return nil
					}

					fmt.Printf("Started %s (pid %d) in the background!\n", activity.Command, activity.Pid)

					return nil
				},
			},
			&cli.Command{
				Name:    "list",
				Aliases: []string{"l", "ls"},
				Usage:   "Lists all running backgrund activities (commands)",
				Action: func(ctx context.Context, c *cli.Command) error {
					err := services.ListActivites()
					if err != nil {
						fmt.Println("Error listing activities (internal error)")
						return nil
					}
					return nil
				},
			},
			&cli.Command{
				Name:    "view",
				Aliases: []string{"v", "show", "log"},
				Usage:   "View live output of process (experimental)",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "pid",
						Usage:   "view process using its pid",
						Aliases: []string{"p"},
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					tmpDir := services.GetTempDirPath()
					var pid int
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
						fmt.Printf("Process \"%s\" not found!\n", c.Args().First())
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

					services.ReadLogFileTail(logFilePath)

					return nil
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
					services.DeleteStoppedActivites()

					if c.Args().Len() < 1 {
						fmt.Println("Please enter a valid pid!")
						return nil
					}

					if c.Bool("pid") == true {
						pid, err := strconv.Atoi(c.Args().First())
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

					name := strings.ToLower(c.Args().First())
					err := services.StopActivityWithName(name, c.Bool("all"))

					if err != nil {
						if err.Error() == "not found" {
							return nil
						}
						fmt.Println("Error while trying to stop process!")
						return nil
					}

					fmt.Println("Stopped process: " + name)
					return nil

				},
			},
			&cli.Command{
				Name:  "temp",
				Usage: "prints temp directory path (for development :D)",
				Action: func(ctx context.Context, c *cli.Command) error {
					fmt.Println(services.GetTempDirPath())
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
