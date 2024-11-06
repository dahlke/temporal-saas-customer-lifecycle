package messages

import (
	"temporal-saas-customer-onboarding/types"

	"go.temporal.io/sdk/workflow"
)

func SetQueryHandlerForState(ctx workflow.Context, state *types.OnboardingWorkflowState) error {
	logger := workflow.GetLogger(ctx)

	err := workflow.SetQueryHandler(ctx, "GetState", func() (types.OnboardingWorkflowState, error) {
		return *state, nil
	})

	if err != nil {
		logger.Error("SetQueryHandler failed for GetState: " + err.Error())
		return err
	}

	return nil
}
