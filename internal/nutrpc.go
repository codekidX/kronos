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
	"go.uber.org/zap"
)

// NutService implements NutServiceServer which is the
// GRPC service governing the state of chrononut
type NutService struct {
	NDB    *NutDatabase
	logger *zap.Logger
}

func (ns *NutService) Init(dbName *string, logger *zap.Logger) {
	db, err := InitializeDB(dbName)
	if err != nil {
		panic(err)
	}
	ns.NDB = db
	ns.logger = logger

	logger.Debug("connected to db",
		zap.Any("stats", db.db.Stats()))

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
			ns.logger.Error("unable to load this task",
				zap.String("name", t.Options.Name),
				zap.String("Ns", t.Options.Ns),
				zap.Error(err))
			continue
		}
	}

	ns.logger.Sugar().Debugf("loaded %d task(s) into engine", len(tasks))
	return nil
}

func (ns *NutService) spawn(opts *proto.TaskOption) error {
	nextTrigger, err := getNextTrigger(opts.GetCronExp())
	if err != nil {
		ns.logger.Error("error while getting trigger",
			zap.String("expr", opts.GetCronExp()),
			zap.String("taskName", opts.Name),
			zap.Error(err))
		return err
	}

	ns.logger.Info("spawning new task",
		zap.String("next", nextTrigger.Format(time.Kitchen)),
		zap.String("taskName", opts.GetName()))
	time.AfterFunc(time.Until(nextTrigger), ns.triggerFunc(opts, time.Now()))
	return nil
}

// Nudge is event in chrononut which takes in a TaskOption with
func (ns *NutService) Nudge(_ context.Context, opts *proto.TaskOption) (*proto.DoneReply, error) {
	if opts.CronExp == "" {
		return nil, fmt.Errorf("schedule not provided %v", opts)
	}

	ns.logger.Debug("inserting task", zap.Any("opts", opts))
	err := ns.NDB.InsertTask(types.Task{Options: opts, Status: types.Active})
	if err != nil {
		ns.logger.Sugar().Error("error while inserting a nudge call",
			zap.Any("opts", opts))
		return nil, err
	}

	// FIXME: best way might be we keep only task id in memory
	// 		  so that it does not keep growing (although we use pointer here)
	err = ns.spawn(opts)
	if err != nil {
		return nil, err
	}
	ns.logger.Sugar().Debug("spawned a new task",
		zap.String("name", opts.Name),
		zap.String("ns", opts.Name))
	return &proto.DoneReply{Ok: true, Message: fmt.Sprintf("Configured : %s", opts.Name)}, nil
}

func (ns *NutService) triggerFunc(opts *proto.TaskOption, start time.Time) func() {
	return func() {
		if !strings.HasPrefix(opts.Url, "http") {
			opts.Url = fmt.Sprintf("http://%s", opts.Url)
			ns.logger.Info("adding http prefix and trying", zap.String("url", opts.Url))
		}
		ns.logger.Debug("nudging..",
			zap.String("url", opts.Url),
			zap.String("payload", string(opts.Data)))
		resp, err := http.Post(opts.Url, "application/json", bytes.NewBuffer(opts.Data))
		if err != nil {
			// TODO: mark task state as errored
			updateerr := ns.NDB.UpdateTaskStatus(opts.Ns, opts.Name, types.Errored)
			if updateerr != nil {
				ns.logger.Error("error while updating task status", zap.Error(updateerr))
			}
			// here create new error artifact
			artinserr := ns.NDB.InsertArtifact(types.TaskArtifact{
				Status:         types.Failure,
				StartTime:      start,
				EndTime:        time.Now(),
				Output:         err.Error(),
				ResponseType:   "None",
				ResponseStatus: 503,
			})
			if artinserr != nil {
				ns.logger.Error("error while inserting failure artifact", zap.Error(artinserr))
			}
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
		artinserr := ns.NDB.InsertArtifact(types.TaskArtifact{
			Status:         types.Success,
			StartTime:      start,
			EndTime:        time.Now(),
			Output:         output,
			ResponseType:   resp.Header.Get("Content-Type"),
			ResponseStatus: resp.StatusCode,
		})
		if artinserr != nil {
			ns.logger.Error("error while inserting success artifact", zap.Error(artinserr))
		}

		ns.spawn(opts)
	}
}

// getNextTrigger gives the next time.Time on which the
// cronExpr should be triggered
func getNextTrigger(cronExpr string) (time.Time, error) {
	expr, err := cronexpr.Parse(cronExpr)
	if err != nil {
		return time.Time{}, fmt.Errorf("not a valid cron expression: %s", cronExpr)
	}

	nextTrigger := expr.Next(time.Now())
	return nextTrigger, nil
}

func (ns *NutService) Cleanup() {
	ns.NDB.cleanup()
}
