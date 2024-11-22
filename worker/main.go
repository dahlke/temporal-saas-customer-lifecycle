package main

import (
	"log"
	"temporal-saas-customer-onboarding/app"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(app.GetClientOptions())
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, app.ONBOARDING_TASK_QUEUE, worker.Options{})
	w.RegisterWorkflow(app.OnboardingWorkflow)
	w.RegisterWorkflow(app.SubscriptionChildWorkflow)
	w.RegisterActivity(app.ChargeCustomer)
	w.RegisterActivity(app.CreateAccount)
	w.RegisterActivity(app.CreateAdminUsers)
	w.RegisterActivity(app.SendClaimCodes)
	w.RegisterActivity(app.SendWelcomeEmail)
	w.RegisterActivity(app.SendFeedbackEmail)
	w.RegisterActivity(app.DeleteAccount)
	w.RegisterActivity(app.DeleteAdminUsers)
	w.RegisterActivity(app.RefundCustomer)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
