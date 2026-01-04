package services

import (
	"fmt"
	"os"
	"runner/internal/models"
	"strings"
	"syscall"
)

func getRunningActivities() ([]models.BackgroundActivity, error) {
	paths, err := readTempDir()
	if err != nil {
		fmt.Println("Error reading activity temp file (internal error)")
		return nil, err
	}
	var activities []models.BackgroundActivity
	for _, p := range paths {
		activity, err := readActivity(p)
		if err != nil {
			fmt.Println("Error reading activity temp file (internal error)")
		}

		activities = append(activities, *activity)
	}

	return activities, nil

}

func StopActivity(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	} else {
		if err := process.Signal(syscall.Signal(0)); err != nil {
			return err
		}
	}

	process.Kill()

	return nil
}

func StopActivityWithName(name string, all bool) error {
	activities, err := getRunningActivities()

	if err != nil {
		fmt.Println("Error getting running background activities (internal error)")
		return err
	}

	// check if process name exists more than once
	count := 0
	var cmd *models.BackgroundActivity
	for _, a := range activities {
		if strings.ToLower(a.Command) == name {
			count += 1
			cmd = &a

			if all {
				pid := cmd.Pid

				if err := StopActivity(pid); err != nil {
					return err
				}
			}
		}
	}

	if all {
		return nil
	}

	if cmd == nil {
		fmt.Printf("Process \"%s\" not found!\n", name)
		return nil
	}

	if count > 1 && !all {
		fmt.Printf("Process with same name (%s) exists more than once!\n", name)
		fmt.Println("Please use the --pid to stop a specific process or --all to quick all processes matching the name!")
		return nil
	}

	pid := cmd.Pid

	if err := StopActivity(pid); err != nil {
		return err
	}

	return nil
}

func ListActivites() error {
	err := DeleteStoppedActivites()
	if err != nil {
		fmt.Println("Error removing already stopped activities (internal error)")
		return err
	}

	activites, err := getRunningActivities()
	if err != nil {
		fmt.Println("Error getting running background activites (internal error)")
		return err
	}
	if len(activites) == 0 {
		fmt.Println("No activities running in the background!")
		return nil
	}
	fmt.Println("All running activities:")
	for _, a := range activites {
		fmt.Printf("\t(%d) %s\n", a.Pid, a.Command)
	}

	return nil
}
