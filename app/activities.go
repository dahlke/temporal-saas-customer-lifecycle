package app

import (
	"context"
	"fmt"

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
	result := fmt.Sprintf("Account for %s was deleted", email)
	return result, nil
}

func CreateAdminUsers(ctx context.Context, emails []string) (string, error) {
	result := fmt.Sprintf("Created admin accounts for %v", emails)
	return result, nil
}

func DeleteAdminUsers(ctx context.Context, emails []string) (string, error) {
	result := fmt.Sprintf("Admin accounts for %v were deleted", emails)
	return result, nil
}

func SendClaimCodes(ctx context.Context, email string, code string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("sending claim code", "email", email, "code", code)
	return "success", nil
}

func SendWelcomeEmail(ctx context.Context, emails []string) (string, error) {
	result := fmt.Sprintf("Sent welcome email to %v", emails)
	return result, nil
}

func SendFeedbackEmail(ctx context.Context, emails []string) (string, error) {
	result := fmt.Sprintf("Sent feedback email to %v", emails)
	return result, nil
}
