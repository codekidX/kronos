package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"nut/gen/proto"
	"nut/pkg/types"
	"strings"
	"time"

	"github.com/gorhill/cronexpr"
)

// NutService implements NutServiceServer which is the
// GRPC service governing the state of chrononut
type NutService struct {
	proto.UnimplementedNutServiceServer
	NDB *NutDatabase
}

func (ns *NutService) Init(dbName *string) {
	db, err := InitializeDB(dbName)
	if err != nil {
		panic(err)
	}
	ns.NDB = db

	err = ns.load()
	if err != nil {
		panic(err)
	}
}

func (ns *NutService) load() error {
	tasks, err := ns.NDB.GetTasks()
	if err != nil {
		return err
	}

	for _, t := range tasks {
		err := ns.spawn(t.Options)
		if err != nil {
			// TODO: log here unable to load the task
			continue
		}
	}

	return nil
}

func (ns *NutService) spawn(opts *proto.TaskOption) error {
	nextTrigger, err := getNextTrigger(opts)
	if err != nil {
		return err
	}

	// TODO: conevrt this to debug log
	fmt.Printf("Spawning %s : next : %s\n", opts.Name, nextTrigger.Format(time.Kitchen))
	time.AfterFunc(time.Until(nextTrigger), ns.triggerFunc(opts, time.Now()))
	return nil
}

// Nudge is event in chrononut which takes in a TaskOption with
func (ns *NutService) Nudge(_ context.Context, opts *proto.TaskOption) (*proto.DoneReply, error) {
	if opts.CronExp == "" {
		return nil, fmt.Errorf("schedule not provided %v", opts)
	}

	err := ns.NDB.InsertTask(types.Task{Options: opts, Status: types.Active})
	if err != nil {
		return nil, err
	}

	err = ns.spawn(opts)
	if err != nil {
		return nil, err
	}
	return &proto.DoneReply{Ok: true, Message: fmt.Sprintf("Configured : %s", opts.Name)}, nil
}

func (ns *NutService) triggerFunc(opts *proto.TaskOption, start time.Time) func() {
	return func() {
		if !strings.HasPrefix(opts.Url, "http") {
			opts.Url = fmt.Sprintf("http://%s", opts.Url)
		}
		resp, err := http.Post(opts.Url, "application/json", bytes.NewBuffer(opts.Data))
		if err != nil {
			// TODO: mark task state as errored
			updateerr := ns.NDB.UpdateTaskStatus(opts.Ns, opts.Name, types.Errored)
			if updateerr != nil {
				fmt.Println("updateerror: ", updateerr.Error())
			}
			// here create new error artifact
			ns.NDB.InsertArtifact(types.TaskArtifact{
				Status:         types.Failure,
				StartTime:      start,
				EndTime:        time.Now(),
				Output:         err.Error(),
				ResponseType:   "None",
				ResponseStatus: 503,
			})
			return
		}

		var output string
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			// TODO: here we log ERROR here
			b = []byte{}
		}
		output = string(b)
		// here create new error artifact
		ns.NDB.InsertArtifact(types.TaskArtifact{
			Status:         types.Success,
			StartTime:      start,
			EndTime:        time.Now(),
			Output:         output,
			ResponseType:   resp.Header.Get("Content-Type"),
			ResponseStatus: resp.StatusCode,
		})

		ns.spawn(opts)
	}
}

func getNextTrigger(opts *proto.TaskOption) (time.Time, error) {
	expr, err := cronexpr.Parse(opts.GetCronExp())
	if err != nil {
		return time.Time{}, fmt.Errorf("not a valid cron expression: %s", opts.GetCronExp())
	}

	nextTrigger := expr.Next(time.Now())
	return nextTrigger, nil
}

func (ns *NutService) Cleanup() {
	ns.NDB.cleanup()
}

func (ns *NutService) mustEmbedUnimplementedNutServiceServer() {
	panic("not implemented") // TODO: Implement
}
