package messages

import (
	"errors"
	"regexp"
	"temporal-saas-customer-onboarding/types"

	"go.temporal.io/sdk/workflow"
)

// "UpdateClaimCode" update handler
func SetUpdateHandlerForAcceptClaimCode(ctx workflow.Context, claimed *bool) (bool, error) {
	logger := workflow.GetLogger(ctx)

	err := workflow.SetUpdateHandlerWithOptions(
		ctx,
		"AcceptClaimCodeUpdate",
		func(ctx workflow.Context, updateInput types.AcceptClaimCodeInput) error {
			*claimed = true
			return nil
		},
		workflow.UpdateHandlerOptions{Validator: validateClaimCode},
	)

	if err != nil {
		logger.Error("SetUpdateHandler failed for UpdateOrder: " + err.Error())
		return false, err
	}

	return *claimed, nil
}

func validateClaimCode(ctx workflow.Context, update types.AcceptClaimCodeInput) error {
	logger := workflow.GetLogger(ctx)
	// Check that the claim code is a 3 letter uppercase string
	re := regexp.MustCompile(`^[A-Z]{3}$`)

	if !re.MatchString(update.ClaimCode) {
		msg := "Rejecting invalid claim code " + update.ClaimCode
		logger.Info(msg)
		return errors.New(msg)
	}

	logger.Info("Updating order, address " + update.ClaimCode)
	return nil
}
