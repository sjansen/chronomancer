package demo

import (
	"context"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

func DemoWorkflow(ctx workflow.Context, emailThreshold time.Duration) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	})
	selector := workflow.NewSelector(ctx)

	emailCtx, cancelEmail := workflow.WithCancel(ctx)
	processingDone := false

	order := workflow.ExecuteActivity(ctx, OrderProcessingActivity)
	selector.AddFuture(order, func(f workflow.Future) {
		cancelEmail()
		processingDone = true
	})

	// send email if processing takes too long
	timer := workflow.NewTimer(emailCtx, emailThreshold)
	selector.AddFuture(timer, func(f workflow.Future) {
		if !processingDone {
			_ = workflow.ExecuteActivity(ctx, SendEmailActivity).Get(ctx, nil)
		}
	})

	// wait for order to be processed, or email sent
	selector.Select(ctx)
	if !processingDone {
		// wait for order processing
		selector.Select(ctx)
	}

	workflow.GetLogger(ctx).Info("✅ Workflow completed.")
	return nil
}

func OrderProcessingActivity(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Info("⛏️ Processing order...")
	timeNeededToProcess := time.Second * time.Duration(rand.Intn(10))
	time.Sleep(timeNeededToProcess)
	logger.Info("😀 Order processed.", "duration", timeNeededToProcess)
	return nil
}

func SendEmailActivity(ctx context.Context) error {
	activity.GetLogger(ctx).Info("✉️  Sending email because of slow processing.")
	return nil
}
