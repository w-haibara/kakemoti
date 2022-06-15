package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/w-haibara/kakemoti/controller/compiler"
	"github.com/w-haibara/kakemoti/db"
	"github.com/w-haibara/kakemoti/pb/worker/workflow"
	"google.golang.org/grpc"
)

func main() {
	go serve()

	name := "sample workflow"
	asl := `
			{
				"StartAt": "Pass State",
				"States": {
				    "Pass State": {
				        "Type": "Pass",
					    "End": true
		            }
				},
				"TimeoutSeconds": 0
			}
			`
	input := `
			{
				"args": ["arg0", "arg1", "arg2"]
			}
			`
	//timeout := 0

	func() {
		fmt.Println("========= compiler.Compile =========")
		fmt.Println("asl:", asl)
		workflow, err := compiler.Compile(context.TODO(), bytes.NewBufferString(asl))
		if err != nil {
			log.Fatal(err)
		}
		db.Save(name, *workflow)
	}()

	time.Sleep(time.Second)

	start(name, input)
}

func serve() {
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

func start(workflowName string, input string) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	c := workflow.NewWorkflowServiceClient(conn)

	response, err := c.Start(context.Background(), &workflow.StartRequest{
		WorkflowName: workflowName,
		Input:        input,
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("output:", response.Output)
}
