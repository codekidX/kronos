package internal

import (
	"context"
	"database/sql"
	"fmt"
	"nut/gen/proto"
	"time"

	"github.com/gorhill/cronexpr"
)

// NutService implements NutServiceServer which is the
// GRPC service governing the state of chrononut
type NutService struct {
	proto.UnimplementedNutServiceServer
	Db *sql.DB
}

// Nudge is event in chrononut which takes in a TaskOption with
func (ns *NutService) Nudge(_ context.Context, opts *proto.TaskOption) (*proto.DoneReply, error) {
	if opts.CronExp == "" {
		return nil, fmt.Errorf("schedule not provided %v", opts)
	}

	nextTrigger, err := getNextTrigger(opts)
	if err != nil {
		return nil, err
	}

	// TODO: here we want to save the nudge to db ..if ns:name is already there then we
	time.AfterFunc(time.Until(nextTrigger), triggerFunc(opts))
	return &proto.DoneReply{Ok: true, Message: fmt.Sprintf("Next at: %s", nextTrigger.Format(time.RFC3339))}, nil
}

func triggerFunc(opts *proto.TaskOption) func() {
	return func() {
		// TODO: here create http request
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

func (ns *NutService) mustEmbedUnimplementedNutServiceServer() {
	panic("not implemented") // TODO: Implement
}
