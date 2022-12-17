package nut

import (
	"context"
	"errors"
	"fmt"
	"nut/gen/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client holds the connection for your crononut instance
// You can build a nut client using `New` function of the nut
// package. `New` expects you to pass a namespace as a parameter
// because a client can belong to only a single namespace.
//
// You can create new clients for different namespaces to avoid
// conflicts between task names.
//
// Example:
//
// 		moneyNut := New("localhost:8999", "moneyService")
// 		adminNut := New("localhost:8999", "adminService")
type Client struct {
	serverAddr string
	conn       *grpc.ClientConn
	ns         string
}

// TaskBuilder helps you build new nut tasks.
// A TaskBuilder can be created using `client.Build` method.
type TaskBuilder struct {
	conn *grpc.ClientConn
	opts *proto.TaskOption
}

func (c *Client) Build(name string) *TaskBuilder {
	opts := &proto.TaskOption{}
	opts.Name = name
	opts.Ns = c.ns
	return &TaskBuilder{opts: opts, conn: c.conn}
}

func (r *TaskBuilder) WithExpression(expression string) *TaskBuilder {
	r.opts.CronExp = expression
	r.opts.IsExact = false
	return r
}

func (r *TaskBuilder) SendPayload(data []byte) *TaskBuilder {
	r.opts.Data = data
	return r
}

func (r *TaskBuilder) Target(url string) *TaskBuilder {
	r.opts.Url = url
	return r
}

func (r *TaskBuilder) At(dateTime time.Time) *TaskBuilder {
	r.opts.CronExp = timeToExpr(dateTime)
	r.opts.IsExact = true
	return r
}

// ForceConnect is a method used to inject RPC connection which can be used for testing
// the nut client.
func (c *Client) ForceConnect(conn *grpc.ClientConn) {
	c.conn = conn
}

// Nudge creates a nudge task in the chrononut server with TaskOption
// that you built using the TaskBuilder.
func (r *TaskBuilder) Nudge() (string, error) {
	if r.opts.Ns == "" {
		return "", errors.New("namespace is required")
	} else if r.opts.Name == "" {
		return "", errors.New("task name is required")
	} else if r.opts.CronExp == "" {
		return "", errors.New("task cannot run without an cron expression")
	} else if r.opts.Url == "" {
		return "", errors.New("cannot be nudged without a target URL")
	}
	rpcClient := proto.NewNutServiceClient(r.conn)
	reply, err := rpcClient.Nudge(context.Background(), r.opts)
	if err != nil {
		return "", err
	}
	// TODO: we need a proper return type instead of stirng
	return reply.Message, nil
}

func timeToExpr(dateTime time.Time) string {
	return fmt.Sprintf("%d %d %d %d * * *", dateTime.Minute(), dateTime.Hour(), dateTime.Day(), dateTime.Month())
}

func New(addr string, ns string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{serverAddr: addr, conn: conn, ns: ns}, nil
}
