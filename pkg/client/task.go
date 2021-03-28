package client

import "time"

type Task struct {
	Uuid        string    `json:"uuid"`
	TaskType    TaskType  `json:"task_type"`
	Input       []byte    `json:"input"`
	Output      []byte    `json:"output"`
	History     []byte    `json:"history"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Error       string    `json:"error"`
	Finish      string    `json:"finish"`
	CreatedTime time.Time `json:"created_time"`
}
