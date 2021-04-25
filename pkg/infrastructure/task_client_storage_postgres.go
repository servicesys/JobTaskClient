package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/satori/go.uuid"
	"github.com/servicesys/JobTaskClient/pkg/client"
)

type TaskClientStoragePostgres struct {
	PoolConn *pgxpool.Pool
}

func NewTaskClientStoragePostgres(pool *pgxpool.Pool) client.TaskClientStorage {

	storagePostgres := &TaskClientStoragePostgres{
		PoolConn: pool,
	}
	return storagePostgres
}

func (t TaskClientStoragePostgres) CreateTaskType(taskType client.TaskType) error {

	queryInsert := `INSERT INTO 
                     job_task.task_type(name, description, input_schema, output_schema, cron_frequent,enable)
                     VALUES($1, $2, $3, $4 , $5, 'S');`

	err := doExecute(t.PoolConn, queryInsert,
		taskType.Name, taskType.Description,
		taskType.InputSchema, taskType.OutputSchema, taskType.CronFrequent)
	return err
}

func (t TaskClientStoragePostgres) GetTaskTypeByName(name string) (client.TaskType, error) {

	fmt.Println(name)
	strQuery := `SELECT name, 
                        description, 
                         input_schema, 
                         output_schema, 
                        cron_frequent 
                 FROM job_task.task_type WHERE name=$1;`
	rows, errQuery := t.PoolConn.Query(context.Background(), strQuery, name)
	mtype := client.TaskType{}

	defer rows.Close()
	if errQuery != nil {
		return mtype, errQuery
	}
	for rows.Next() {
		rows.Scan(&mtype.Name, &mtype.Description, &mtype.InputSchema, &mtype.OutputSchema, &mtype.CronFrequent)
	}
	return mtype, nil
}

func (t TaskClientStoragePostgres) GetAllTaskNotStartedByType(name string) ([]client.Task, error) {

	strQuery :=
		`SELECT uuid, task_type_name , input , output , start_time, end_time, error, finish, created_time,
         tt.description, tt.input_schema, tt.output_schema, tt.cron_frequent
         FROM job_task.task t 
         INNER JOIN job_task.task_type tt ON (t.task_type_name=tt.name)
         WHERE (finish is null OR finish ='N') AND (enable='S' OR  enable IS NULL) AND task_type_name = $1;`

	tasks := make([]client.Task, 0)
	rows, errQuery := t.PoolConn.Query(context.Background(), strQuery, name)
	task := client.Task{}

	defer rows.Close()
	if errQuery != nil {
		return tasks, errQuery
	}
	for rows.Next() {
		task = client.Task{}
		task.TaskType = client.TaskType{}
		startTime := sql.NullTime{}
		endTime := sql.NullTime{}
		strError := sql.NullString{}
		strFinish := sql.NullString{}

		errorScan := rows.Scan(&task.Uuid,
			&task.TaskType.Name,
			&task.Input,
			&task.Output,
			&startTime,
			&endTime,
			&strError,
			&strFinish,
			&task.CreatedTime,
			&task.TaskType.Description,
			&task.TaskType.InputSchema,
			&task.TaskType.OutputSchema,
			&task.TaskType.CronFrequent,
		)

		if errorScan != nil {
			return tasks, errorScan
		}
		task.StartTime = startTime.Time
		task.EndTime = endTime.Time
		task.Error = strError.String
		task.Finish = strFinish.String
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (t TaskClientStoragePostgres) AddTask(task client.Task) error {

	task.Uuid = uuid.NewV4().String()
	query := ` INSERT INTO job_task.task (uuid, task_type_name, input, created_time) VALUES($1 , $2, $3 ,now());`
	err := doExecute(t.PoolConn, query, task.Uuid, task.TaskType.Name, task.Input)
	return err
}

func (t TaskClientStoragePostgres) UpdateTask(task client.Task) error {

	query := `UPDATE job_task.task 
              SET output=$1 , history= $2, start_time=$3, end_time= $4 , error=$5, finish=$6 WHERE uuid=$7;`

	err := doExecute(t.PoolConn, query, task.Output, task.History, task.StartTime, task.EndTime,
		task.Error, task.Finish, task.Uuid)
	return err
}

func doExecute(dbPool *pgxpool.Pool, query string, args ...interface{}) error {

	_, err := dbPool.Exec(context.Background(), query, args...)

	return err
}
