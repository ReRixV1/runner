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

func GetTempDirPath() string {
	osTempDir := os.TempDir()
	tempDir := filepath.Join(osTempDir, "net.rerix.runner")
	return tempDir
}

func EnsureTempDirectory() error {
	tempDir := GetTempDirPath()
	//fmt.Println("TempDir: " + tempDir)
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error ensuring temp directory (internal error)")
		return err
	}
	return nil
}

func getLogFileName(processName string, suffix int) string {
	tmpDir := GetTempDirPath()

	name := processName
	if suffix != 0 {
		name += "-" + strconv.Itoa(suffix)
	}
	name += ".log"
	if _, err := os.Stat(filepath.Join(tmpDir, name)); err == nil {
		return getLogFileName(processName, suffix+1)
	}

	return name
}

func createAndOpenLogFile(name string) (*os.File, error) {
	tempDir := GetTempDirPath()
	path := filepath.Join(tempDir, name)
	return os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	//_, err := os.Create(path)

	//if err != nil {
	//	fmt.Println("Error trying to create temp log file (internal error)")
	//	return nil, nil
	//}

	//return os.Open(path)

}

func readTempDir() ([]string, error) {
	tempDir := GetTempDirPath()
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
	tempDir := GetTempDirPath()

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

func GetActivity(path string) (*models.BackgroundActivity, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var activity *models.BackgroundActivity

	err = json.Unmarshal(f, &activity)

	return activity, err
}

func removeFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func DeleteStoppedActivites() error {
	activities, err := getRunningActivities()
	if err != nil {
		fmt.Println("Error removing already stopped activities (internal error)")
		return err
	}

	tempDir := GetTempDirPath()

	for _, a := range activities {
		pid := a.Pid
		process, err := os.FindProcess(pid)

		if err != nil {
			path := filepath.Join(tempDir, strconv.Itoa(pid)+".json")
			stopProcess(path)
		} else {
			if err := process.Signal(syscall.Signal(0)); err != nil {
				path := filepath.Join(tempDir, strconv.Itoa(pid)+".json")
				stopProcess(path)
			}
		}
	}

	return nil
}
