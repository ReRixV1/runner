package services

import (
	"errors"
	"fmt"
	"os"
	"runner/internal/models"
	"strings"
	"syscall"
)

func getRunningActivities() ([]models.BackgroundActivity, error) {
	paths, err := readTempDir()
	if err != nil {
		fmt.Println("Error getting temp files (internal error)")
		return nil, err
	}
	var activities []models.BackgroundActivity
	for _, p := range paths {
		if !strings.HasSuffix(p, ".json") {
			continue
		}
		activity, err := GetActivity(p)
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

func GetPids(name string) ([]*int, error) {
	activities, err := getRunningActivities()

	if err != nil {
		fmt.Println("Error getting running background activities (internal error)")
		return nil, err
	}

	var pids []*int
	for _, a := range activities {
		if strings.ToLower(a.Command) == name {
			pids = append(pids, &a.Pid)
		}
	}

	return pids, nil
}

func StopActivityWithName(name string, all bool) error {
	pids, _ := GetPids(name)

	if all {
		for _, p := range pids {
			if err := StopActivity(*p); err != nil {
				return err
			}
		}
		return nil
	}

	if len(pids) == 0 {
		fmt.Printf("Process \"%s\" not found\n", name)
		return errors.New("not found")
	}

	if len(pids) > 1 && !all {
		fmt.Printf("Process with same name (%s) exists more than once!\n", name)
		fmt.Println("Please use the --pid to stop a specific process or --all to quick all processes matching the name!")
		return nil
	}

	pid := *pids[0]

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
