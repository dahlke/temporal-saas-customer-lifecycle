package app

import (
	"math/rand"
	"time"
)

const (
	LIFECYCLE_TASK_QUEUE             = "LIFECYCLE_TASK_QUEUE"
	ACCEPTANCE_TIME                  = 120 // Time in seconds to wait for claim codes to be accepted
	SCENARIO_HAPPY_PATH              = "HAPPY_PATH"
	SCENARIO_FLAKEY_API              = "FLAKEY_API"
	SCENARIO_RECOVERABLE_FAILURE     = "RECOVERABLE_FAILURE"
	SCENARIO_NON_RECOVERABLE_FAILURE = "NON_RECOVERABLE_FAILURE"
	SCENARIO_CHILD_WORKFLOW          = "CHILD_WORKFLOW"
)

func generateNewClaimCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 3)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
