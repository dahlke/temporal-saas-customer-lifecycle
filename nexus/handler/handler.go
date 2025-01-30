package handler

import (
	"context"
	"fmt"
	"temporal-saas-customer-lifecycle/app"
	"temporal-saas-customer-lifecycle/types"

	"go.temporal.io/sdk/client"

	"github.com/google/uuid"
	"github.com/nexus-rpc/sdk-go/nexus"
	"go.temporal.io/sdk/temporalnexus"
)

var BillingOperation = temporalnexus.NewWorkflowRunOperation(
	app.BillingOperationName,
	app.SubscriptionBillingWorkflow,
	func(ctx context.Context, input types.LifecycleInput, soo nexus.StartOperationOptions) (client.StartWorkflowOptions, error) {
		return client.StartWorkflowOptions{ID: fmt.Sprintf("subscription-billing-%v-%v", input.AccountName, uuid.New().String())}, nil
	},
)
