package app

import (
	"fmt"
	"temporal-saas-customer-onboarding/messages"
	"temporal-saas-customer-onboarding/types"
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	ACCEPTANCE_TIME = 120
)

func OnboardingWorkflow(ctx workflow.Context, input types.OnboardingWorkflowInput) (string, error) {
	// logger := workflow.GetLogger(ctx)

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	state := types.OnboardingWorkflowState{
		AccountName: input.AccountName,
		Emails:      input.Emails,
		ClaimCodes:  make([]types.ClaimCodeStatus, len(input.Emails)),
	}

	// Initialize claim codes for each email
	claimCodes := []string{"XXX", "YYY"}
	for i, email := range input.Emails {
		state.ClaimCodes[i] = types.ClaimCodeStatus{
			Email:     email,
			Code:      claimCodes[i],
			IsClaimed: false,
		}
	}

	err := messages.SetQueryHandlerForState(ctx, &state)
	if err != nil {
		return "", err
	}

	// TODO: throw and catch errors
	// TODO: create a compensation class, do compensations
	// TODO: custom search attributes
	// TODO: re-send claim codes signal
	// TODO: add logging of the results
	// TODO: resend welcome email signal
	// TODO: resend claim codes signal
	// TODO: comments
	// TODO: set IsClaimed based on which claim code is accepted

	var chargeResult string
	err = workflow.ExecuteActivity(ctx, ChargeCustomer, input.AccountName).Get(ctx, &chargeResult)

	if err != nil {
		return "", err
	}

	var createAccountResult string
	err = workflow.ExecuteActivity(ctx, CreateAccount, input.AccountName).Get(ctx, &createAccountResult)

	if err != nil {
		return "", err
	}

	var createAdminUsersResult string
	err = workflow.ExecuteActivity(ctx, CreateAdminUsers, input.Emails).Get(ctx, &createAdminUsersResult)

	if err != nil {
		return "", err
	}

	// Make the claim code a hash of the emails?
	for _, claimCode := range state.ClaimCodes {
		var sendClaimCodeResult string
		err = workflow.ExecuteActivity(ctx, SendClaimCodes, input.AccountName, claimCode.Code).Get(ctx, &sendClaimCodeResult)
		if err != nil {
			return "", err
		}
	}

	if err != nil {
		return "", err
	}

	// Create a pointer to track the claimed status
	var claimed bool
	claimed, err = messages.SetUpdateHandlerForAcceptClaimCode(ctx, &claimed)
	if err != nil {
		return "", err
	}

	// Wait for up to ACCEPTANCE_TIME seconds for the update
	ok, _ := workflow.AwaitWithTimeout(ctx, time.Second*ACCEPTANCE_TIME, func() bool {
		return claimed
	})

	// If the update wasn't received or was false, fail the workflow
	if !ok {
		return "", fmt.Errorf("claim codes not accepted within %d seconds", ACCEPTANCE_TIME)
	}

	var sendWelcomeEmailResult string
	err = workflow.ExecuteActivity(ctx, SendWelcomeEmail, input.Emails).Get(ctx, &sendWelcomeEmailResult)

	if err != nil {
		return "", err
	}

	var sendFeedbackEmailResult string
	err = workflow.ExecuteActivity(ctx, SendFeedbackEmail, input.Emails).Get(ctx, &sendFeedbackEmailResult)

	if err != nil {
		return "", err
	}

	return sendFeedbackEmailResult, nil
}
