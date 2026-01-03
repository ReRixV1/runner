package models

type BackgroundActivity struct {
	Command   string   `json:"command"`
	Pid       int      `json:"pid"`
	Arguments []string `json:"arguments"`
}
