package main

import (
	"context"
	"fmt"
	"log"
	"temporal-saas-customer-lifecycle/app"
	"temporal-saas-customer-lifecycle/types"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(app.GetClientOptions())
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "lifecycle-workflow-" + uuid.New().String(),
		TaskQueue: app.LIFECYCLE_TASK_QUEUE,
	}

	// Start the Workflow
	accountName := "Temporal"
	emails := []string{"neil@dahlke.io", "neil.dahlke@temporal.io"}
	price := 10.0
	input := types.LifecycleWorkflowInput{
		AccountName: accountName,
		Emails:      emails,
		Price:       price,
		Scenario:    app.SCENARIO_CHILD_WORKFLOW,
	}

	wf, err := c.ExecuteWorkflow(
		context.Background(),
		options,
		app.LifecycleWorkflow,
		input,
	)

	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// Get the results
	var result string
	err = wf.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}

	printResults(result, wf.GetID(), wf.GetRunID())
}

func printResults(result string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s\n\n", result)
}
