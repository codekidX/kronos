package main

import (
	"net"
	"net/http"
	"nut/gen/proto"
	"nut/internal"
	"nut/internal/admin"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// initialize chrononut config

	// initialize db

	// initialize Task channels
	// initialize Log channels
	lis, err := net.Listen("tcp", "localhost:8122")
	if err != nil {
		panic(err)
	}

	// zap logger init w.r.t environment
	logger := internal.CreateLogger()

	// initialize GRPC server
	var opts []grpc.ServerOption
	rpcServer := grpc.NewServer(opts...)
	nutService := &internal.NutService{}
	dao := nutService.Init(nil, logger)
	proto.RegisterNutServiceServer(rpcServer, nutService)
	logger.Info("Starting Nut RPC Server", zap.String("port", "8122"))

	if err := admin.ServePackaged(); err != nil {
		// start GRPC server
		rpcServer.Serve(lis)
		logger.Error("Error cannot serve packaged", zap.Any("ref", err))
	} else {
		// we also add the Admin APIs now
		admin.AddRoutes(dao)
		// start GRPC server
		go rpcServer.Serve(lis)
		logger.Info("Starting Nut Admin Panel", zap.String("port", "8121"))
		http.ListenAndServe("localhost:8121", nil)
	}
}
