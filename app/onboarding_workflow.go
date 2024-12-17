package app

import (
	"fmt"
	"temporal-saas-customer-lifecycle/messages"
	"temporal-saas-customer-lifecycle/types"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var lifecycleStatusKey = temporal.NewSearchAttributeKeyKeyword("LifecycleStatus")

// LifecycleWorkflow orchestrates the lifecycle process for a new customer.
func LifecycleWorkflow(ctx workflow.Context, input types.LifecycleWorkflowInput) (string, error) {
	logger := workflow.GetLogger(ctx)

	// Log the start of the workflow
	logger.Info("Lifecycle workflow started", "account_name", input.AccountName, "emails", input.Emails)

	// Initialize workflow state
	state := types.LifecycleWorkflowState{
		AccountName: input.AccountName,
		Emails:      input.Emails,
		Price:       input.Price,
		ClaimCodes:  make([]types.ClaimCodeStatus, len(input.Emails)),
		Progress:    0,
		Status:      "UNINITIALIZED",
	}

	// Initialize claim codes for each email
	claimCodes := []string{generateNewClaimCode(), generateNewClaimCode()}
	for i, email := range input.Emails {
		state.ClaimCodes[i] = types.ClaimCodeStatus{
			Email:     email,
			Code:      claimCodes[i],
			IsClaimed: false,
		}
	}

	logger.Info("Claim codes after assignment", "claimCodes", state.ClaimCodes)

	// Set initial search attribute for the workflow
	workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet("STARTED"))
	state.Progress = 10
	state.Status = "STARTED"

	// Define retry policy for activities
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Second * 10,
		MaximumAttempts:    10,
		/*
			// TODO
			NonRetryableErrorTypes: []string{
			},
		*/
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

	// Set query handler for the workflow state
	err = messages.SetQueryHandlerForState(ctx, &state)
	if err != nil {
		return "", err
	}

	// Update search attribute to indicate charging phase
	workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet("CHARGING"))
	state.Progress = 20
	state.Status = "CHARGING"

	// Charge customer
	var chargeResult string
	err = workflow.ExecuteActivity(ctx, ChargeCustomer, input).Get(ctx, &chargeResult)
	if err != nil {
		return "", err
	}
	saga.AddCompensation(RefundCustomer, input)
	logger.Info("Successfully charged customer", "result", chargeResult)

	// Update search attribute to indicate account creation phase
	workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet("CREATING_ACCOUNT"))
	state.Progress = 30
	state.Status = "CREATING_ACCOUNT"

	// Create account
	var createAccountResult string
	err = workflow.ExecuteActivity(ctx, CreateAccount, input).Get(ctx, &createAccountResult)
	if err != nil {
		return "", err
	}
	saga.AddCompensation(DeleteAccount, input)
	logger.Info("Successfully created account", "result", createAccountResult)

	// Update search attribute to indicate admin user creation phase
	workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet("CREATING_ADMIN_USERS"))
	state.Progress = 40
	state.Status = "CREATING_ADMIN_USERS"

	// Create admin users
	var createAdminUsersResult string
	err = workflow.ExecuteActivity(ctx, CreateAdminUsers, input).Get(ctx, &createAdminUsersResult)
	if err != nil {
		return "", err
	}
	saga.AddCompensation(DeleteAdminUsers, input)
	logger.Info("Successfully created admin users", "result", createAdminUsersResult)

	if input.Scenario == SCENARIO_UNEXPECTED_BUG {
		// Simulate bug
		// NOTE: comment out this line to see the happy path
		panic("Simulated bug - fix me!")
	}

	// Update search attribute to indicate claim code sending phase
	workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet("SENDING_CLAIM_CODES"))
	state.Progress = 50
	state.Status = "SENDING_CLAIM_CODES"

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
	workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet("WAITING_FOR_CLAIM_CODES"))
	state.Progress = 60
	state.Status = "WAITING_FOR_CLAIM_CODES"

	// Await signal message to update address
	logger.Info(fmt.Sprintf("Waiting up to %d seconds for claim codes", ACCEPTANCE_TIME))
	var signal messages.ResendClaimCodesSignal

	// Get signal channel for resend claim codes
	claimCodesSignalChan := messages.GetSignalChannelForResendClaimCodes(ctx)

	// Goroutine to handle resend claim codes signal
	workflow.Go(ctx, func(ctx workflow.Context) {
		for {
			selector := workflow.NewSelector(ctx)
			selector.AddReceive(claimCodesSignalChan, func(c workflow.ReceiveChannel, more bool) {
				c.Receive(ctx, &signal)
				logger.Info("Received resend claim codes signal")

				// Generate new claim codes for each email
				for i := range state.ClaimCodes {
					state.ClaimCodes[i].Code = generateNewClaimCode()
				}

				logger.Info("New claim codes after assignment", "claimCodes", state.ClaimCodes)

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
		workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet("CODE_NOT_CLAIMED"))
		state.Status = "CODE_NOT_CLAIMED"
		logger.Info("Claim codes not accepted within %d seconds", ACCEPTANCE_TIME)
	}

	if state.Status != "CODE_NOT_CLAIMED" {
		// Update the claim status in the workflow state
		for i := range state.ClaimCodes {
			if state.ClaimCodes[i].Code == acceptedCode {
				state.ClaimCodes[i].IsClaimed = true
				break
			}
		}

		workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet("SENDING_WELCOME_EMAIL"))
		state.Progress = 70
		state.Status = "SENDING_WELCOME_EMAIL"

		// Send welcome email
		var sendWelcomeEmailResult string
		err = workflow.ExecuteActivity(ctx, SendWelcomeEmail, input).Get(ctx, &sendWelcomeEmailResult)
		if err != nil {
			logger.Error("Failed to send welcome email", "error", err)
			return "", err
		}
		logger.Info("Successfully sent welcome email", "result", sendWelcomeEmailResult)

		// Clear saga compensations as the group is onboarded
		saga.ClearCompensations()

		// Wait before sending feedback email
		logger.Info("Waiting 1 seconds before sending feedback email")
		workflow.Sleep(ctx, time.Second*1)

		workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet("SENDING_FEEDBACK_EMAIL"))
		state.Progress = 80
		state.Status = "SENDING_FEEDBACK_EMAIL"

		// Send feedback email
		var sendFeedbackEmailResult string
		err = workflow.ExecuteActivity(ctx, SendFeedbackEmail, input).Get(ctx, &sendFeedbackEmailResult)
		if err != nil {
			logger.Error("Failed to send feedback email", "error", err)
			return "", err
		}
		logger.Info("Successfully sent feedback email", "result", sendFeedbackEmailResult)

		workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet("ONBOARDED"))
		state.Progress = 90
		state.Status = "ONBOARDED"

		if input.Scenario == SCENARIO_CHILD_WORKFLOW {
			// Start the subscription child workflow
			ChildWorkflowOptions := workflow.ChildWorkflowOptions{
				WorkflowID:        fmt.Sprintf("subscription-%v-%v", input.AccountName, uuid.New().String()),
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			}
			ctx = workflow.WithChildOptions(ctx, ChildWorkflowOptions)
			err := workflow.ExecuteChildWorkflow(ctx, SubscriptionChildWorkflow, input).Get(ctx, nil)
			if err != nil {
				logger.Error("Failed to start subscription child workflow", "error", err)
				return "", err
			}
			logger.Info("Started Child Workflow: " + ChildWorkflowOptions.WorkflowID)
		} else {
			// Create a channel to receive the cancel subscription signal
			cancelSubscriptionSignalChan := messages.GetSignalChannelForCancelSubscription(ctx)

			subscriptionCanceled := false
			numRenews := 0
			for {
				logger.Info("Waiting for 3 seconds to charge the customer or until a cancel subscription signal is received")
				// Wait for 3 seconds or until a cancel subscription signal is received
				selector := workflow.NewSelector(ctx)
				selector.AddReceive(cancelSubscriptionSignalChan, func(c workflow.ReceiveChannel, more bool) {
					// Break the loop when the signal is received
					logger.Info("Received cancel subscription signal")
					subscriptionCanceled = true
					workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet("SUBSCRIPTION_CANCELED"))
					state.Status = "SUBSCRIPTION_CANCELED"
				})
				selector.AddFuture(workflow.NewTimer(ctx, time.Second*3), func(f workflow.Future) {
					// Timer expired, continue the loop
				})
				selector.Select(ctx)

				// Check if the subscription was canceled
				if subscriptionCanceled {
					break
				}

				// Execute the charge activity
				var chargeResult string
				err = workflow.ExecuteActivity(ctx, ChargeCustomer, input).Get(ctx, &chargeResult)
				if err != nil {
					logger.Error("Failed to charge customer", "error", err)
					return "", err
				}

				logger.Info("Successfully charged customer", "result", chargeResult)

				numRenews++
				renewStatus := fmt.Sprintf("RENEWED_%d", numRenews)
				workflow.UpsertTypedSearchAttributes(ctx, lifecycleStatusKey.ValueSet(renewStatus))
				state.Status = renewStatus
			}
		}

		// TODO: we could also clean up the admin users and account here
		// saga.AddCompensation(DeleteAccount, input)
		// saga.AddCompensation(DeleteAdminUsers, input)
		state.Progress = 100
	}

	return state.Status, nil
}
