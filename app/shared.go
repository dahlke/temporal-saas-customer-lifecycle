package app

import (
	"math/rand"
	"time"
)

const (
	LIFECYCLE_TASK_QUEUE             = "subscription-lifecycle-task-queue"
	ACCEPTANCE_TIME                  = 120 // Time in seconds to wait for claim codes to be accepted
	SCENARIO_HAPPY_PATH              = "HAPPY_PATH"
	SCENARIO_FLAKEY_API              = "FLAKEY_API"
	SCENARIO_RECOVERABLE_FAILURE     = "RECOVERABLE_FAILURE"
	SCENARIO_NON_RECOVERABLE_FAILURE = "NON_RECOVERABLE_FAILURE"
	SCENARIO_CHILD_WORKFLOW          = "CHILD_WORKFLOW"
	SCENARIO_NEXUS_WORKFLOW          = "NEXUS"
)

const BillingServiceName = "billing-service"
const BillingOperationName = "charge-customer"

type BillingInput struct {
	AccountName string
	Price       float64
	Scenario    string
}

type BillingOutput struct {
	Message string
}

func GenerateNewClaimCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 3)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
