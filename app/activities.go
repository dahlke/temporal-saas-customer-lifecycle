package app

import (
	"context"

	"go.temporal.io/sdk/activity"
)

func ChargeCustomer(ctx context.Context, customerID string, amount float64) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("charging customer", "customer_id", customerID, "amount", amount)
	return "success", nil
}

func RefundCustomer(ctx context.Context, customerID string, amount float64) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("refunding customer", "customer_id", customerID, "amount", amount)
	return "success", nil
}

func CreateAccount(ctx context.Context, email string, password string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("creating account", "email", email)
	return "success", nil
}

func DeleteAccount(ctx context.Context, email string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("deleting account", "email", email)
	return "success", nil
}

func CreateAdminUsers(ctx context.Context, emails []string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("creating admin users", "emails", emails)
	return "success", nil
}

func DeleteAdminUsers(ctx context.Context, emails []string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("deleting admin users", "emails", emails)
	return "success", nil
}

func SendClaimCodes(ctx context.Context, email string, code string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("sending claim code", "email", email, "code", code)
	return "success", nil
}

func SendWelcomeEmail(ctx context.Context, emails []string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("sending welcome email", "emails", emails)
	return "success", nil
}

func SendFeedbackEmail(ctx context.Context, emails []string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("sending feedback email", "emails", emails)
	return "success", nil
}
