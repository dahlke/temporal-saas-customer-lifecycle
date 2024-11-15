package messages

import (
	"errors"
	"fmt"
	"temporal-saas-customer-onboarding/types"

	"go.temporal.io/sdk/workflow"
)

// "UpdateClaimCode" update handler
func SetUpdateHandlerForAcceptClaimCode(ctx workflow.Context, claimed *bool, acceptedCode *string, state *types.OnboardingWorkflowState) (bool, error) {
	logger := workflow.GetLogger(ctx)

	// Set an update handler with options for the "AcceptClaimCodeUpdate" operation
	err := workflow.SetUpdateHandlerWithOptions(
		ctx,
		"AcceptClaimCodeUpdate",
		func(ctx workflow.Context, updateInput types.AcceptClaimCodeInput) error {
			// Only set claimed if the code is valid
			if isValidClaimCode(updateInput.ClaimCode, state) {
				*claimed = true
				*acceptedCode = updateInput.ClaimCode
				return nil
			}
			// Return an error if the claim code is not found
			return fmt.Errorf("claim code %s not found in workflow state", updateInput.ClaimCode)
		},
		workflow.UpdateHandlerOptions{
			// Validator function to validate the claim code
			Validator: func(ctx workflow.Context, input types.AcceptClaimCodeInput) error {
				return validateClaimCode(ctx, input, state)
			},
		},
	)

	// Log an error if setting the update handler fails
	if err != nil {
		logger.Error("SetUpdateHandler failed for UpdateOrder: " + err.Error())
		return false, err
	}

	// Return the claimed status
	return *claimed, nil
}

// Function to validate the claim code
func validateClaimCode(ctx workflow.Context, update types.AcceptClaimCodeInput, state *types.OnboardingWorkflowState) error {
	logger := workflow.GetLogger(ctx)

	// Check if the code exists in the workflow state
	if !isValidClaimCode(update.ClaimCode, state) {
		msg := "Rejecting unknown claim code: " + update.ClaimCode
		logger.Info(msg)
		return errors.New(msg)
	}

	// Check if the code has already been claimed
	for _, code := range state.ClaimCodes {
		if code.Code == update.ClaimCode && code.IsClaimed {
			msg := "Rejecting already claimed code: " + update.ClaimCode
			logger.Info(msg)
			return errors.New(msg)
		}
	}

	// Log that a valid claim code was received
	logger.Info("Valid claim code received: " + update.ClaimCode)
	return nil
}

// Helper function to check if a claim code exists in the workflow state
func isValidClaimCode(code string, state *types.OnboardingWorkflowState) bool {
	for _, claimCode := range state.ClaimCodes {
		if claimCode.Code == code {
			return true
		}
	}
	return false
}
