package app

import (
	"context"
	"errors"
	"temporal-saas-customer-lifecycle/types"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

func ChargeCustomer(ctx context.Context, input types.LifecycleWorkflowInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("charging customer", "customer_id", input.AccountName, "amount", input.Price)
	time.Sleep(1 * time.Second)

	if input.Scenario == SCENARIO_NON_RECOVERABLE_FAILURE {
		return "", temporal.NewNonRetryableApplicationError(
			"charge customer activity failed, card invalid", "activityFailure",
			errors.New("charge customer API failure, card invalid"),
		)
	}

	return "success", nil
}

func RefundCustomer(ctx context.Context, input types.LifecycleWorkflowInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("refunding customer", "customer_id", input.AccountName, "amount", input.Price)
	time.Sleep(1 * time.Second)
	return "success", nil
}

func CreateAccount(ctx context.Context, input types.LifecycleWorkflowInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("creating account", "email", input.AccountName)
	time.Sleep(1 * time.Second)
	return "success", nil
}

func DeleteAccount(ctx context.Context, input types.LifecycleWorkflowInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("deleting account", "email", input.AccountName)
	time.Sleep(1 * time.Second)
	return "success", nil
}

func CreateAdminUsers(ctx context.Context, input types.LifecycleWorkflowInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("creating admin users", "emails", input.Emails)
	time.Sleep(1 * time.Second)

	if input.Scenario == SCENARIO_FLAKEY_API {
		info := activity.GetInfo(ctx)
		if info.Attempt < 5 {
			return "failure", errors.New("create admin users activity failed, API unavailable")
		}
	}

	return "success", nil
}

func DeleteAdminUsers(ctx context.Context, input types.LifecycleWorkflowInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("deleting admin users", "emails", input.Emails)
	time.Sleep(1 * time.Second)
	return "success", nil
}

func SendClaimCodes(ctx context.Context, input types.LifecycleWorkflowInput, claimCode string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("sending claim code", "email", input.AccountName, "code", claimCode)
	time.Sleep(1 * time.Second)
	return "success", nil
}

func SendWelcomeEmail(ctx context.Context, input types.LifecycleWorkflowInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("sending welcome email", "emails", input.Emails)
	time.Sleep(1 * time.Second)
	return "success", nil
}

func SendFeedbackEmail(ctx context.Context, input types.LifecycleWorkflowInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("sending feedback email", "emails", input.Emails)
	time.Sleep(1 * time.Second)
	return "success", nil
}
