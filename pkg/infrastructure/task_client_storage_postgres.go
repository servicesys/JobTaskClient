package infrastructure

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
)

type TaskClientStoragePostgres struct {
	dbConnection *pgx.Conn
}

func NewTaskClientStoragePostgres() client.TaskClientStorage {

	/**
	DB.HOST="localhost"
	DB.PORT="5432"
	DB.USER="valter"
	DB.PASS="valter"
	DB.NAME="app_sistema"
	*/
	connection := Connect("localhost", 5432, "valter", "valter", "app_sistema")
	storagePostgres := &TaskClientStoragePostgres{
		dbConnection: connection,
	}

	return storagePostgres
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

func (t TaskClientStoragePostgres) CreateTaskType(taskType server.TaskType) error {

	queryInsert := `INSERT INTO 
                     job_task.task_type(name, description, input_schema, output_schema, cron_frequent)
                     VALUES($1, $2, $3, $4 , $5);`

	err := doExecute(t.dbConnection, queryInsert,
		taskType.Name, taskType.Description,
		taskType.InputSchema, taskType.OutputSchema, taskType.CronFrequent)
	return err
}

func (t TaskClientStoragePostgres) GetTaskTypeByName(name string) (server.TaskType, error) {

	strQuery := `SELECT name, description, input_schema, output_schema, cron_frequent FROM job_task.task_type WHERE name=$1;`
	rows, errQuery := t.dbConnection.Query(context.Background(), strQuery, name)
	mtype := server.TaskType{}

	defer rows.Close()
	if errQuery != nil {
		return mtype, errQuery
	}
	for rows.Next() {
		rows.Scan(&mtype.Name, &mtype.Description, &mtype.InputSchema, &mtype.OutputSchema, &mtype.CronFrequent)
	}
	return mtype, nil
}

func doExecute(db *pgx.Conn, query string, args ...interface{}) error {

	_, err := db.Exec(context.Background(), query, args...)

	return err
}
