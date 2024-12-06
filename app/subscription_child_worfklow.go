package app

import (
	"fmt"
	"temporal-saas-customer-lifecycle/messages"
	"temporal-saas-customer-lifecycle/types"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func SubscriptionChildWorkflow(ctx workflow.Context, input types.LifecycleWorkflowInput) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Subscription child workflow started", "accountName", input.AccountName)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    1 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    30 * time.Second,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Create a channel to receive the cancel subscription signal
	cancelSubscriptionSignalChan := messages.GetSignalChannelForCancelSubscription(ctx)

	subscriptionCanceled := false
	numRenews := 0
	for {
		logger.Info("Waiting for 10 seconds to charge the customer or until a cancel subscription signal is received")
		// Wait for 10 seconds or until a cancel subscription signal is received
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(cancelSubscriptionSignalChan, func(c workflow.ReceiveChannel, more bool) {
			// Break the loop when the signal is received
			logger.Info("Received cancel subscription signal")
			subscriptionCanceled = true
			workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet("SUBSCRIPTION_CANCELED"))
		})
		selector.AddFuture(workflow.NewTimer(ctx, time.Second*10), func(f workflow.Future) {
			// Timer expired, continue the loop
		})
		selector.Select(ctx)

		// Check if the subscription was canceled
		if subscriptionCanceled {
			break
		}

		// Execute the charge activity
		var chargeResult string
		err := workflow.ExecuteActivity(ctx, ChargeCustomer, input).Get(ctx, &chargeResult)
		if err != nil {
			logger.Error("Failed to charge customer", "error", err)
			return "", err
		}

		numRenews++
		workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet(fmt.Sprintf("RENEWED_%d", numRenews)))
		logger.Info("Successfully charged customer", "result", chargeResult)
	}

	return "", nil
}
