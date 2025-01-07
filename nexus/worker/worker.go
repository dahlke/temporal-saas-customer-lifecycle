package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"temporal-saas-customer-lifecycle/app"
	"temporal-saas-customer-lifecycle/nexus/handler"

	"github.com/nexus-rpc/sdk-go/nexus"
)

func main() {
	c, err := client.Dial(app.GetClientOptions(true))
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, app.GetEnv("NEXUS_BILLING_TASK_QUEUE", "subscription-billing-task-queue"), worker.Options{})
	// TODO: take this from the constants in shared.go?
	service := nexus.NewService(app.BillingServiceName)
	err = service.Register(handler.BillingOperation)
	if err != nil {
		log.Fatalln("Unable to register operations", err)
	}
	w.RegisterNexusService(service)
	w.RegisterWorkflow(app.SubscriptionBillingWorkflow)
	w.RegisterActivity(app.ChargeCustomer)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
