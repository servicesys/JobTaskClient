package client

type TaskClientStorage interface {
	CreateTaskType(taskType TaskType) error
	GetTaskTypeByName(name string) (TaskType, error)
	GetAllTaskNotStartedByType(name string) ([]Task, error)
	AddTask(task Task) error
	UpdateTask(task Task) error
}
