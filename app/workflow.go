package app

import (
	"fmt"
	"temporal-saas-customer-onboarding/messages"
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	ACCEPTANCE_TIME = 120
)

func OnboardingWorkflow(ctx workflow.Context, name string) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	// TODO: codec server
	// TODO: create a compensation class, do compensations
	// TODO: custom search attributes
	// TODO: throw and catch errors
	// TODO: re-send claim codes signal
	// TODO: add logging

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
	var emails []string = []string{"neil@dahlke.io", "neil.dahlke@temporal.io"}
	err = workflow.ExecuteActivity(ctx, CreateAdminUsers, emails).Get(ctx, &createAdminUsersResult)

	if err != nil {
		return "", err
	}

	var createSupportChannelResult string
	err = workflow.ExecuteActivity(ctx, CreateSupportChannel, name).Get(ctx, &createSupportChannelResult)

	if err != nil {
		return "", err
	}

	var sendClaimCodesResult string
	var claimCodes []string = []string{"XXX", "YYY"}
	err = workflow.ExecuteActivity(ctx, SendClaimCodes, name, claimCodes).Get(ctx, &sendClaimCodesResult)

	if err != nil {
		return "", err
	}

	/*
		// Signal to receive the claim code
		c := messages.GetSignalChannelForAcceptClaimCode(ctx)
		claimCodeStatus, _ := c.ReceiveWithTimeout(ctx, time.Second*ACCEPTANCE_TIME, nil)

		// If the signal was not received within the timeout, fail the workflow
		if !claimCodeStatus {
			return "", fmt.Errorf("claim codes not accepted within %d seconds", ACCEPTANCE_TIME)
		}
	*/

	// Update to receive the claim code
	var claimCodeStatus bool
	claimCodeStatus, err = messages.SetUpdateHandlerForAcceptClaimCode(ctx)

	if err != nil {
		return "", err
	}

	workflow.AwaitWithTimeout(ctx, time.Minute, func() bool {
		return claimCodeStatus
	})

	// If the update or signal was not received within the timeout, fail the workflow
	if !claimCodeStatus {
		return "", fmt.Errorf("claim codes not accepted within %d seconds", ACCEPTANCE_TIME)
	}

	if err != nil {
		return "", err
	}

	// TODO: re-send welcome codes signal
	// TODO: query out welcome email sent
	var sendWelcomeEmailResult string
	err = workflow.ExecuteActivity(ctx, SendWelcomeEmail, emails).Get(ctx, &sendWelcomeEmailResult)

	if err != nil {
		return "", err
	}

	var sendFeedbackEmailResult string
	err = workflow.ExecuteActivity(ctx, SendFeedbackEmail, emails).Get(ctx, &sendFeedbackEmailResult)

	if err != nil {
		return "", err
	}

	return sendFeedbackEmailResult, nil
}
