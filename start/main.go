package main

import (
	"context"
	"fmt"
	"log"
	"temporal-saas-customer-onboarding/app"
	"temporal-saas-customer-onboarding/types"

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
	accountName := "Temporal"
	emails := []string{"neil@dahlke.io", "neil.dahlke@temporal.io"}
	price := 10.0
	input := types.OnboardingWorkflowInput{AccountName: accountName, Emails: emails, Price: price}

	wf, err := c.ExecuteWorkflow(
		context.Background(),
		options,
		app.OnboardingWorkflow,
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
