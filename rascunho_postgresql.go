package main

import (
	"JobTask/pkg/infrastructure"
	"fmt"
)

func main() {

	fmt.Println("DB")
	taskClientStorage := infrastructure.NewTaskClientStoragePostgres()

	/*
	taskType := server.TaskType{
		Name:         "OLA",
		Description:  "Ola 123333 ",
		InputSchema:  nil,
		OutputSchema: nil,
		CronFrequent: "@every 2s",
		TaskJobRef:   nil,
	}
	taskClientStorage.CreateTaskType(taskType)*/

	taskT , error := taskClientStorage.GetTaskTypeByName("OLA")
	fmt.Println(taskT.Description)
	fmt.Println(error)

}
