package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runner/internal/models"
	"strconv"
	"syscall"
)

func getTempDirPath() string {
	osTempDir := os.TempDir()
	tempDir := filepath.Join(osTempDir, "net.rerix.runner")
	return tempDir
}

func EnsureTempDirectory() error {
	tempDir := getTempDirPath()
	fmt.Println("TempDir: " + tempDir)
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error ensuring temp directory (internal error)")
		return err
	}
	return nil
}

func readTempDir() ([]string, error) {
	tempDir := getTempDirPath()
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		fmt.Println("Error reading temp directory (internal error)")
		return nil, err
	}
	var files []string
	for _, e := range entries {
		name := e.Name()
		path := filepath.Join(tempDir, name)
		files = append(files, path)
	}
	return files, nil
}

func WriteActivity(activity models.BackgroundActivity) error {
	tempDir := getTempDirPath()

	data, err := json.Marshal(activity)
	if err != nil {
		return err
	}

	path := filepath.Join(tempDir, strconv.Itoa(activity.Pid)+".json")
	err = os.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func readActivity(path string) (*models.BackgroundActivity, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var activity *models.BackgroundActivity

	err = json.Unmarshal(f, &activity)

	return activity, err
}

func removeActivity(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func deleteStoppedActivites() error {
	activities, err := getRunningActivities()
	if err != nil {
		return err
	}

	tempDir := getTempDirPath()

	for _, a := range activities {
		pid := a.Pid
		process, err := os.FindProcess(pid)

		if err != nil {
			path := filepath.Join(tempDir, strconv.Itoa(pid)+".json")
			if err = removeActivity(path); err != nil {
				return err
			}
		} else {
			if err := process.Signal(syscall.Signal(0)); err != nil {
				path := filepath.Join(tempDir, strconv.Itoa(pid)+".json")
				if err = removeActivity(path); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
