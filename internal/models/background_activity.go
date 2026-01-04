package models

type BackgroundActivity struct {
	Command   string   `json:"command"`
	Pid       int      `json:"pid"`
	LogFile   string   `json:"logFile"`
	Arguments []string `json:"arguments"`
}
