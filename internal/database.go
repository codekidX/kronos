package internal

import (
	"database/sql"
	"nut/gen/proto"
	"nut/pkg/types"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Persistance interface {
	GetTasks() ([]types.Task, error)
	InsertTask(types.Task) error
	GetArtifacts() ([]types.TaskArtifact, error)
	InsertArtifact(types.TaskArtifact) error
}

type NutDatabase struct {
	db *sql.DB
}

func (nd *NutDatabase) cleanup() {
	_, err := nd.db.Exec(`DELETE FROM tasks`)
	if err != nil {
		panic(err)
	}
	_, err = nd.db.Exec(`DELETE FROM artifacts`)
	if err != nil {
		panic(err)
	}
	_, err = nd.db.Exec(`VACUUM`)
	if err != nil {
		panic(err)
	}
}

func (nd *NutDatabase) UpdateTaskStatus(ns, name string, status types.TaskStatus) error {
	stmt, err := nd.db.Prepare(`UPDATE tasks SET status = ? WHERE ns = ? AND name = ?`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(status, ns, name)
	if err != nil {
		return err
	}
	return nil
}

// INTERFACE IMPLEMENTATION

func (nd *NutDatabase) GetTasks() ([]types.Task, error) {
	var tasks []types.Task
	rows, err := nd.db.Query("SELECT * from tasks WHERE status != 3")
	if err != nil {
		return tasks, err
	}

	for rows.Next() {
		t := types.Task{
			Options: &proto.TaskOption{},
		}
		err = rows.Scan(&t.ID, &t.Options.Ns, &t.Options.Name, &t.Type, &t.Options.Data, &t.Options.Url, &t.Options.CronExp)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (nd *NutDatabase) InsertTask(task types.Task) error {
	stmt, err := nd.db.Prepare(
		`INSERT INTO tasks (
            ns,
            name,
            status,
            data,
            url,
            cron_exp) VALUES(?,?,?,?,?,?)`)

	if err != nil {
		return err
	}

	// TODO: maybe use last inserted id !
	_, err = stmt.Exec(task.Options.Ns, task.Options.Name, task.Status, string(task.Options.Data), task.Options.Url, task.Options.CronExp)
	if err != nil {
		return err
	}

	return nil
}

func (nd *NutDatabase) GetArtifacts() ([]types.TaskArtifact, error) {
	panic("not implemented") // TODO: Implement
}

func (nd *NutDatabase) InsertArtifact(a types.TaskArtifact) error {
	stmt, err := nd.db.Prepare(
		`INSERT INTO artifacts (output,
			status,
			response_type,
			start_dtm,
			end_dtm,
			response_status) VALUES (?,?,?,?,?,?)`)
	if err != nil {
		// TODO: log warning error here
		return err
	}

	_, err = stmt.Exec(a.Output, a.Status, a.ResponseType, a.StartTime.Format(time.Kitchen), a.EndTime.Format(time.Kitchen), a.ResponseStatus)
	if err != nil {
		return err
	}

	return nil
}

func InitializeDB(name *string) (*NutDatabase, error) {
	if name == nil {
		tmp := "nut.db"
		name = &tmp
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", filepath.Join(homeDir, *name))
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS tasks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            ns VARCHAR NOT NULL,
            name VARCHAR NOT NULL,
            status VARCHAR NOT NULL,
            data VARCHAR,
            url VARCHAR NOT NULL,
            cron_exp VARCHAR,
            UNIQUE(ns, name))`)

	if err != nil {
		return nil, err
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS artifacts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            output VARCHAR,
            status INTEGER  NOT NULL,
            response_type VARCHAR,
            start_dtm DATETIME NOT NULL,
            end_dtm DATETIME NOT NULL, 
            response_status VARCHAR NOT NULL)`)
	if err != nil {
		return nil, err
	}

	return &NutDatabase{db: db}, nil
}
