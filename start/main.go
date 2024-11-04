package main

import (
	"context"
	"fmt"
	"log"
	"temporal-saas-customer-onboarding/app"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func main() {

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "onboarding-workflow-" + uuid.New().String(),
		TaskQueue: app.OnboardingTaskQueue,
	}

	// Start the Workflow
	name := "World"
	we, err := c.ExecuteWorkflow(context.Background(), options, app.OnboardingWorkflow, name)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// Get the results
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}

	printResults(result, we.GetID(), we.GetRunID())
}

func printResults(result string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s\n\n", result)
}
