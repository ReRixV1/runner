package services

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runner/internal/models"
	"syscall"
)

func ExecCommand(commands ...string) {
	cmd := exec.Command(commands[0], commands[1:]...)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("Error executing command")
		fmt.Println(err)
		return
	}

	fmt.Println(string(stdout))
}
func StartProcessInBackground(commands ...string) (*models.BackgroundActivity, error) {
	cmd := exec.Command(commands[0], commands[1:]...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	logName := getLogFileName(commands[0], 0)
	fout, err := createAndOpenLogFile(logName)
	if err != nil {
		fmt.Println("Error opening log file (internal error)")
		return nil, nil
	}

	devNull, err := os.OpenFile("/dev/null", os.O_RDONLY, 0)
	if err != nil {
		fmt.Println("Error opening /dev/null (internal error)")
		return nil, err
	}

	cmd.Stdin = devNull
	cmd.Stdout = fout
	cmd.Stderr = fout

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	activity := &models.BackgroundActivity{
		Command:   commands[0],
		Pid:       cmd.Process.Pid,
		LogFile:   logName,
		Arguments: commands[1:],
	}

	WriteActivity(*activity)
	cmd.Process.Release()

	return activity, nil
}

func stopProcess(path string) error {
	tmpDir := GetTempDirPath()
	var activity *models.BackgroundActivity

	f, err := os.ReadFile(path)
	err = json.Unmarshal(f, &activity)
	if err != nil {
		return err
	}
	logPath := filepath.Join(tmpDir, activity.LogFile)

	if err = removeFile(logPath); err != nil {
		return err
	}

	if err = removeFile(path); err != nil {
		return err
	}
	return nil
}

func restartProcess(pid int) error {
	return nil
}
