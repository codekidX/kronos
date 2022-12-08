package types

import (
	"nut/gen/proto"
	"time"
)

type TaskStatus int
type TaskType int
type ArtifactStatus int

const (
	Active TaskStatus = iota
	Stopped
	Finished
	Errored

	Trigger TaskType = iota
	Nudge
	Dated

	Success ArtifactStatus = iota
	Failure
)

type Task struct {
	ID      int
	Status  TaskStatus
	Type    TaskType
	Options *proto.TaskOption
}

type TaskArtifact struct {
	ResponseStatus int
	Output         string
	Status         ArtifactStatus
	ResponseType   string
	StartTime      time.Time
	EndTime        time.Time
}
