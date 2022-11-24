package types

import (
	"nut/gen/proto"
)

type TaskStatus int
type TaskType int
type ArtifactStatus int

const (
	Active TaskStatus = iota
	Stopped
	Finished

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
	Output        string
	Status        ArtifactStatus
	Response_type string
}
