package main

import (
	"log"

	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"

	"github.com/sjansen/chronomancer/internal/demo"
)

func main() {
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
		Logger:   logur.LoggerToKV(logrusadapter.New(logrus.New())),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "demo", worker.Options{
		MaxConcurrentActivityExecutionSize: 3,
	})

	w.RegisterWorkflow(demo.DemoWorkflow)

	w.RegisterActivity(demo.OrderProcessingActivity)
	w.RegisterActivity(demo.SendEmailActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
