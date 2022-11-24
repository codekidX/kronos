package nut

import (
	"context"
	"errors"
	"nut/gen/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	serverAddr string
	conn       *grpc.ClientConn
	ns         string
}

type Request struct {
	conn *grpc.ClientConn
	opts *proto.TaskOption
}

func (c *Client) Build(namespace string, name string) *Request {
	opts := &proto.TaskOption{}
	opts.Name = name
	opts.Ns = namespace
	return &Request{opts: opts, conn: c.conn}
}

func (r *Request) WithExpression(expression string, isExact bool) *Request {
	r.opts.CronExp = expression
	r.opts.IsExact = isExact
	return r
}

func (r *Request) SendPayload(data []byte) *Request {
	r.opts.Data = data
	return r
}

func (r *Request) Target(url string) *Request {
	r.opts.Url = url
	return r
}

// ForceConnect is a method used to inject RPC connection which can be used for testing
// the nut client.
func (c *Client) ForceConnect(conn *grpc.ClientConn) {
	c.conn = conn
}

func (r *Request) Nudge() (string, error) {
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
	return reply.Message, nil
}

func New(addr string, ns string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{serverAddr: addr, conn: conn, ns: ns}, nil
}
