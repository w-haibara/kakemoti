package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/w-haibara/kakemoti/compiler"
	"github.com/w-haibara/kakemoti/worker"
)

func main() {
	//name := "sample workflow"
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

	ctx := context.TODO()

	fmt.Println("========= compiler.Compile =========")
	fmt.Println("asl:", asl)
	workflow, err := compiler.Compile(ctx, bytes.NewBufferString(asl))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========= workflow.Exec =========")
	fmt.Println("input:", input)
	res, err := worker.Exec(ctx, nil, *workflow, bytes.NewBufferString(input))
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("result:", string(res))
}
