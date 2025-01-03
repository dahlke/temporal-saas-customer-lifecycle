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
	c, err := client.Dial(app.GetClientOptions())
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, app.GetEnv("TEMPORAL_NEXUS_BILLING_TASK_QUEUE", "billing"), worker.Options{})
	service := nexus.NewService(app.NEXUS_BILLING_SERVICE_NAME)
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
