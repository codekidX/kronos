package internal

import (
	"database/sql"
	"nut/pkg/types"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Persistance interface {
	GetTasks() (types.Task, error)
	InsertTask(types.Task) error
	GetArtifacts() (types.TaskArtifact, error)
	InsertArtifact(types.TaskArtifact) error
}

type NutDatabase struct{}

func (nd *NutDatabase) GetTasks() (types.Task, error) {
	panic("not implemented") // TODO: Implement
}

func (nd *NutDatabase) InsertTask(_ types.Task) error {
	panic("not implemented") // TODO: Implement
}

func (nd *NutDatabase) GetArtifacts() (types.TaskArtifact, error) {
	panic("not implemented") // TODO: Implement
}

func (nd *NutDatabase) InsertArtifact(_ types.TaskArtifact) error {
	panic("not implemented") // TODO: Implement
}

func InitializeDB(name *string) (*sql.DB, error) {
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
            typ VARCHAR NOT NULL,
            data VARCHAR,
            url VARCHAR NOT NULL,
            cron_exp VARCHAR,
            trigger_date VARCHAR,
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

	return db, nil
}
