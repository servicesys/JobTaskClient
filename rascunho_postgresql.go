package main

import (
	"JobTaskClient/pkg/client"
	"JobTaskClient/pkg/infrastructure"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"time"
)

func main() {

	fmt.Println("DB")
	connection := Connect("localhost", 5432, "valter", "valter", "app_sistema")
	taskClientStorage := infrastructure.NewTaskClientStoragePostgres(connection)

	tasks, errorGetTask := taskClientStorage.GetAllTaskNotStartedByType("HELLO")
	fmt.Println(errorGetTask)
	fmt.Println(len(tasks))
	//fmt.Println(tasks)
	//testes(taskClientStorage)
	//createTaskType(taskClientStorage)
	createTask(taskClientStorage)

}

func createTaskType(taskClientStorage client.TaskClientStorage) {
	helloTaskType := client.TaskType{
		Name:         "HELLO",
		Description:  "Hello for users",
		InputSchema:  nil,
		OutputSchema: nil,
		CronFrequent: "@every 5s",
	}

	worldTaskType := client.TaskType{
		Name:         "WORLD",
		Description:  "World for users",
		InputSchema:  nil,
		OutputSchema: nil,
		CronFrequent: "@every 15s",
	}

	taskClientStorage.CreateTaskType(helloTaskType)
	taskClientStorage.CreateTaskType(worldTaskType)
}

func createTask(storage client.TaskClientStorage) {

	textoJSon := ` { "title" : " Hello world task job input" , "text" :  "HELLO"}`

	input := []byte(textoJSon)

	listHelloTask := make([]client.Task, 2)

	for i := 0; i < 2; i++ {

		listHelloTask[i] = client.Task{
			//Uuid: "hello" + string(i),
			TaskType: client.TaskType{
				Name:         "HELLO",
			},
			Input:       input,
		}
		errorAddTask := storage.AddTask(listHelloTask[i])
		fmt.Println(errorAddTask)

	}

}

func testes(taskClientStorage client.TaskClientStorage) {
	taskType := client.TaskType{
		Name:         "OLA",
		Description:  "Ola 123333 ",
		InputSchema:  getSchema(),
		OutputSchema: getSchema(),
		CronFrequent: "@every 2s",
	}
	errorCreate := taskClientStorage.CreateTaskType(taskType)
	fmt.Println(errorCreate)
	taskT, error := taskClientStorage.GetTaskTypeByName("OLA")
	fmt.Println(taskT)
	fmt.Println(error)

	textoJSon := ` { "title" : " Hello world task job input" , "text" :  "HELLO"} `

	input := []byte(textoJSon)
	taskOla := client.Task{
		Uuid: "hello" + string(123),
		TaskType: client.TaskType{
			Name:         "OLA",
			Description:  "",
			InputSchema:  nil,
			OutputSchema: nil,
			CronFrequent: "",
		},
		Input:       input,
		Output:      nil,
		History:     nil,
		StartTime:   time.Time{},
		EndTime:     time.Time{},
		Error:       "",
		Finish:      "",
		CreatedTime: time.Time{},
	}

	errorAddTask := taskClientStorage.AddTask(taskOla)
	fmt.Println(errorAddTask)

	taskOla.Output = []byte(textoJSon)
	taskOla.StartTime = time.Now()
	taskOla.EndTime = time.Now()
	taskOla.Finish = "S"
	taskOla.Uuid = "497d7778-9cc1-49ca-b85d-16707a77f5ae"
	errorUpdateTask := taskClientStorage.UpdateTask(taskOla)
	fmt.Println(errorUpdateTask)

	tasks, errorGetTask := taskClientStorage.GetAllTaskNotStartedByType("OLA")
	fmt.Println(errorGetTask)
	fmt.Println(len(tasks))
}


func Connect(host string, port int, user string, pass string, db string) *pgx.Conn {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, pass, db)

	dbConn, err := pgx.Connect(context.Background(), psqlInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		panic(err)
	}
	return dbConn

}

func getSchema() []byte {

	SCHEMA := `{
     "$id": "https://qri.io/schema/",
    "$comment" : "sample comment",
    "title": "Texto Blog",
    "type": "object",
	"properties": {
		"title": {

			"title": "Titulo",
			"type": "string",
			"default": "",
			"examples": [
				"Este e um texto de exemplo"
			],
			"pattern": "^.*$"
		},
		"text": {

			"title": "Texto",
			"type": "string",
			"default": "",
			"examples": [
				"<p>Este e  o corpot do texto texto de exemplo</p>"
			],
			"pattern": "^.*$"
		}
	},
	"required": [
		"title",
		"text"
	]
}
`
	return []byte(SCHEMA)
}


