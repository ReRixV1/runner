package services

import (
	"fmt"
	"runner/internal/models"
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

func ListActivites() error {
	err := deleteStoppedActivites()
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
