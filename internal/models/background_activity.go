package models

import "bytes"

type BackgroundActivity struct {
	cmdStr  string
	outputb *bytes.Buffer
}
