package app

import (
	"math/rand"
	"time"
)

const (
	LIFECYCLE_TASK_QUEUE    = "LIFECYCLE_TASK_QUEUE"
	ACCEPTANCE_TIME         = 120 // Time in seconds to wait for claim codes to be accepted
	SCENARIO_HAPPY_PATH     = "SCENARIO_HAPPY_PATH"
	SCENARIO_FLAKEY_API     = "SCENARIO_FLAKEY_API"
	SCENARIO_UNEXPECTED_BUG = "SCENARIO_UNEXPECTED_BUG"
	SCENARIO_EXPECTED_ERROR = "SCENARIO_EXPECTED_ERROR"
	SCENARIO_CHILD_WORKFLOW = "SCENARIO_CHILD_WORKFLOW"
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
