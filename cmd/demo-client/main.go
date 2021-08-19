package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"go.temporal.io/sdk/client"

	"github.com/sjansen/chronomancer/internal/demo"
)

func main() {
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	now := time.Now()
	// TODO use sync.Pool to reuse entropy source
	entropy := ulid.Monotonic(rand.New(rand.NewSource(now.UnixNano())), 0)

	ulid, err := ulid.New(ulid.Timestamp(now), entropy)
	if err != nil {
		log.Fatalln("Unable to create workflow ID", err)
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        "demo_" + ulid.String(),
		TaskQueue: "demo",
	}

	ctx := context.TODO()
	we, err := c.ExecuteWorkflow(ctx, workflowOptions, demo.DemoWorkflow, time.Second*3)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
