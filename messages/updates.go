package messages

import (
	"errors"
	"regexp"

	"go.temporal.io/sdk/workflow"
)

// "UpdateClaimCode" update handler
func SetUpdateHandlerForAcceptClaimCode(ctx workflow.Context) (bool, error) {
	logger := workflow.GetLogger(ctx)

	var codeAcceptedStatus bool

	err := workflow.SetUpdateHandlerWithOptions(
		ctx,
		"AcceptClaimCode",
		func(ctx workflow.Context, updateInput AcceptClaimCodeInput) (bool, error) {
			codeAcceptedStatus = true
			return codeAcceptedStatus, nil
		},
		workflow.UpdateHandlerOptions{Validator: validateClaimCode},
	)

	if err != nil {
		logger.Error("SetUpdateHandler failed for UpdateOrder: " + err.Error())
		return false, err
	}

	return codeAcceptedStatus, nil
}

func validateClaimCode(ctx workflow.Context, update AcceptClaimCodeInput) error {
	logger := workflow.GetLogger(ctx)

	re := regexp.MustCompile(`^[A-Z]{3}$`)

	if !re.MatchString(update.ClaimCode) {
		msg := "Rejecting invalid claim code " + update.ClaimCode
		logger.Info(msg)
		return errors.New(msg)
	}

	logger.Info("Updating order, address " + update.ClaimCode)
	return nil
}
