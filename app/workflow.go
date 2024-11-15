package app

import (
	"fmt"
	"temporal-saas-customer-onboarding/messages"
	"temporal-saas-customer-onboarding/types"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	ACCEPTANCE_TIME = 120
)

var onboardingStatusKey = temporal.NewSearchAttributeKeyKeyword("OnboardingStatus")

// OnboardingWorkflow orchestrates the onboarding process for a new customer.
func OnboardingWorkflow(ctx workflow.Context, input types.OnboardingWorkflowInput) (string, error) {
	logger := workflow.GetLogger(ctx)

	// Log the start of the workflow
	logger.Info("Onboarding workflow started", "account_name", input.AccountName, "emails", input.Emails)

	// Set initial search attribute for the workflow
	workflow.UpsertTypedSearchAttributes(ctx, onboardingStatusKey.ValueSet("STARTED"))

	// Define retry policy for activities
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Second * 10,
		MaximumAttempts:    10,
	}

	// Set activity options with retry policy
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 60,
		RetryPolicy:         retrypolicy,
	}

	// Apply activity options to the context
	ctx = workflow.WithActivityOptions(ctx, options)
	var err error
	var saga Saga

	// Ensure compensation in case of error
	defer func() {
		if err != nil {
			disconnectedCtx, _ := workflow.NewDisconnectedContext(ctx)
			saga.Compensate(disconnectedCtx)
		}
	}()

	// Initialize workflow state
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

	// Set query handler for the workflow state
	err = messages.SetQueryHandlerForState(ctx, &state)
	if err != nil {
		return "", err
	}

	// Update search attribute to indicate charging phase
	workflow.UpsertTypedSearchAttributes(ctx, onboardingStatusKey.ValueSet("CHARGING"))

	// Charge customer
	var chargeResult string
	err = workflow.ExecuteActivity(ctx, ChargeCustomer, input).Get(ctx, &chargeResult)
	if err != nil {
		return "", err
	}
	saga.AddCompensation(RefundCustomer, input)
	logger.Info("Successfully charged customer", "result", chargeResult)

	// Update search attribute to indicate account creation phase
	workflow.UpsertTypedSearchAttributes(ctx, onboardingStatusKey.ValueSet("CREATING_ACCOUNT"))

	// Create account
	var createAccountResult string
	err = workflow.ExecuteActivity(ctx, CreateAccount, input).Get(ctx, &createAccountResult)
	if err != nil {
		return "", err
	}
	saga.AddCompensation(DeleteAccount, input)
	logger.Info("Successfully created account", "result", createAccountResult)

	// Update search attribute to indicate admin user creation phase
	workflow.UpsertTypedSearchAttributes(ctx, onboardingStatusKey.ValueSet("CREATING_ADMIN_USERS"))

	// Create admin users
	var createAdminUsersResult string
	err = workflow.ExecuteActivity(ctx, CreateAdminUsers, input).Get(ctx, &createAdminUsersResult)
	if err != nil {
		return "", err
	}
	saga.AddCompensation(DeleteAdminUsers, input)
	logger.Info("Successfully created admin users", "result", createAdminUsersResult)

	// Simulate bug
	// panic("Simulated bug - fix me!")

	// Update search attribute to indicate claim code sending phase
	workflow.UpsertTypedSearchAttributes(ctx, onboardingStatusKey.ValueSet("SENDING_CLAIM_CODES"))

	// Send claim codes
	for _, claimCode := range state.ClaimCodes {
		var sendClaimCodeResult string
		err = workflow.ExecuteActivity(ctx, SendClaimCodes, input, claimCode.Code).Get(ctx, &sendClaimCodeResult)
		if err != nil {
			logger.Error("Failed to send claim code", "error", err, "email", claimCode.Email)
			return "", err
		}
		logger.Info("Successfully sent claim code", "result", sendClaimCodeResult, "email", claimCode.Email)
	}

	// Update search attribute to indicate waiting for claim codes
	workflow.UpsertTypedSearchAttributes(ctx, onboardingStatusKey.ValueSet("WAITING_FOR_CLAIM_CODES"))

	// Await signal message to update address
	logger.Info("Waiting up to 60 seconds for resend claim codes")
	var signal messages.ResendClaimCodesSignal

	// Get signal channel for resend claim codes
	signalChan := messages.GetSignalChannelForResendClaimCodes(ctx)

	// Goroutine to handle resend claim codes signal
	workflow.Go(ctx, func(ctx workflow.Context) {
		for {
			selector := workflow.NewSelector(ctx)
			selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
				c.Receive(ctx, &signal)
				logger.Info("Received resend claim codes signal")

				// Resend claim codes for each email
				for _, claimCode := range state.ClaimCodes {
					var sendClaimCodeResult string
					err := workflow.ExecuteActivity(ctx, SendClaimCodes, input, claimCode.Code).Get(ctx, &sendClaimCodeResult)
					if err != nil {
						logger.Error("Failed to resend claim code", "error", err, "email", claimCode.Email)
						continue
					}
					logger.Info("Successfully resent claim code", "result", sendClaimCodeResult, "email", claimCode.Email)
				}
			})
			selector.Select(ctx)
		}
	})

	// Create a pointer to track the claimed status and the accepted code
	var claimed bool
	var acceptedCode string
	claimed, err = messages.SetUpdateHandlerForAcceptClaimCode(ctx, &claimed, &acceptedCode, &state)
	if err != nil {
		return "", err
	}

	// Wait for up to ACCEPTANCE_TIME seconds for the update
	ok, _ := workflow.AwaitWithTimeout(ctx, time.Second*ACCEPTANCE_TIME, func() bool {
		return claimed
	})

	// If the update wasn't received or was false, fail the workflow
	if !ok {
		workflow.UpsertTypedSearchAttributes(ctx, onboardingStatusKey.ValueSet("CODE_NOT_CLAIMED"))
		return "", fmt.Errorf("claim codes not accepted within %d seconds", ACCEPTANCE_TIME)
	}

	// Update the claim status in the workflow state
	for i := range state.ClaimCodes {
		if state.ClaimCodes[i].Code == acceptedCode {
			state.ClaimCodes[i].IsClaimed = true
			break
		}
	}

	workflow.UpsertTypedSearchAttributes(ctx, onboardingStatusKey.ValueSet("SENDING_WELCOME_EMAIL"))

	var sendWelcomeEmailResult string
	err = workflow.ExecuteActivity(ctx, SendWelcomeEmail, input).Get(ctx, &sendWelcomeEmailResult)
	if err != nil {
		logger.Error("Failed to send welcome email", "error", err)
		return "", err
	}
	logger.Info("Successfully sent welcome email", "result", sendWelcomeEmailResult)

	logger.Info("Waiting 10 seconds before sending feedback email")
	workflow.Sleep(ctx, time.Second*10)

	workflow.UpsertTypedSearchAttributes(ctx, onboardingStatusKey.ValueSet("SENDING_FEEDBACK_EMAIL"))

	var sendFeedbackEmailResult string
	err = workflow.ExecuteActivity(ctx, SendFeedbackEmail, input).Get(ctx, &sendFeedbackEmailResult)
	if err != nil {
		logger.Error("Failed to send feedback email", "error", err)
		return "", err
	}
	logger.Info("Successfully sent feedback email", "result", sendFeedbackEmailResult)

	workflow.UpsertTypedSearchAttributes(ctx, onboardingStatusKey.ValueSet("COMPLETED"))

	/*
		// Now we can wait a period of time and charge the customer on a recurring basis.
		// You will want to clear out the saga compensations and arguments before you do this.

		// TODO: do this as a detached child workflow to allow this workflow to complete?

		for {
			// Wait for 30 seconds
			workflow.Sleep(ctx, time.Second*30)

			// Execute the charge activity
			var chargeResult string
			err = workflow.ExecuteActivity(ctx, ChargeCustomer, input).Get(ctx, &chargeResult)
			if err != nil {
				logger.Error("Failed to charge customer", "error", err)
				return "", err
			}
			logger.Info("Successfully charged customer", "result", chargeResult)
		}
	*/

	return sendFeedbackEmailResult, nil
}
