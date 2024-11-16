package app

import "go.temporal.io/sdk/workflow"

type Saga struct {
	compensations []any   // List of compensation activities to be executed
	arguments     [][]any // Arguments for each compensation activity
}

// AddCompensation adds a compensation activity and its parameters to the saga
func (s *Saga) AddCompensation(activity any, parameters ...any) {
	s.compensations = append(s.compensations, activity) // Append the activity to the compensations list
	s.arguments = append(s.arguments, parameters)       // Append the parameters to the arguments list
}

// ClearCompensations clears the compensations and arguments from the saga so we can start fresh.
func (s *Saga) ClearCompensations() {
	s.compensations = nil
	s.arguments = nil
}

// Compensate executes all compensation activities in reverse order
func (s Saga) Compensate(ctx workflow.Context) {
	logger := workflow.GetLogger(ctx)         // Get a logger from the workflow context
	logger.Info("Saga compensations started") // Log the start of compensation

	// Compensate in the reverse order that activities were applied.
	for i := len(s.compensations) - 1; i >= 0; i-- {
		// Execute the compensation activity with its arguments
		err := workflow.ExecuteActivity(ctx, s.compensations[i], s.arguments[i]...).Get(ctx, nil)
		if err != nil {
			// Log an error if the compensation activity fails
			logger.Error("Executing compensation failed", "Error", err)
		}
	}
}
