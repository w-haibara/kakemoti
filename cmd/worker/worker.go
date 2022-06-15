package main

import (
	"net"

	"github.com/w-haibara/kakemoti/pb/worker/workflow"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err.Error())
	}

	s := workflow.Server{}

	grpcServer := grpc.NewServer()

	workflow.RegisterWorkflowServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err.Error())
	}
}
