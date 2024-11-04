package app

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func OnboardingWorkflow(ctx workflow.Context, name string) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	// TODO: compensations
	// TODO: custom search attributes

	var chargeResult string
	err := workflow.ExecuteActivity(ctx, ChargeCustomer, name).Get(ctx, &chargeResult)

	if err != nil {
		return "", err
	}

	var createAccountResult string
	err = workflow.ExecuteActivity(ctx, CreateAccount, name).Get(ctx, &createAccountResult)

	if err != nil {
		return "", err
	}

	var createAdminUsersResult string
	err = workflow.ExecuteActivity(ctx, CreateAdminUsers, name).Get(ctx, &createAdminUsersResult)

	if err != nil {
		return "", err
	}

	var createSupportChannelResult string
	err = workflow.ExecuteActivity(ctx, CreateSupportChannel, name).Get(ctx, &createSupportChannelResult)

	if err != nil {
		return "", err
	}

	// TODO: re-send claim codes signal
	// TODO: Signal claim code accepter
	var sendClaimCodesResult string
	err = workflow.ExecuteActivity(ctx, SendClaimCodes, name).Get(ctx, &sendClaimCodesResult)

	if err != nil {
		return "", err
	}

	// TODO: re-send welcome codes signal
	// TODO: query out welcome email sent
	var sendWelcomeEmailResult string
	err = workflow.ExecuteActivity(ctx, SendWelcomeEmail, name).Get(ctx, &sendWelcomeEmailResult)

	if err != nil {
		return "", err
	}

	// TODO: query out feedback email
	var sendFeedbackEmailResult string
	err = workflow.ExecuteActivity(ctx, SendFeedbackEmail, name).Get(ctx, &sendFeedbackEmailResult)

	if err != nil {
		return "", err
	}

	return sendFeedbackEmailResult, nil
}
