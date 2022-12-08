package main

import (
	"fmt"
	"net"
	"net/http"
	"nut/gen/proto"
	"nut/internal"

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

	// initialize GRPC server
	var opts []grpc.ServerOption
	rpcServer := grpc.NewServer(opts...)
	nutService := &internal.NutService{}
	nutService.Init(nil)
	proto.RegisterNutServiceServer(rpcServer, nutService)
	fmt.Println("Starting Nut RPC Server")
	// start GRPC server
	go rpcServer.Serve(lis)

	// create HTTP server for serving admin panel
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Awesome!!")
	})
	fmt.Println("Starting Nut Admin Panel")
	http.ListenAndServe("localhost:8121", nil)
}
