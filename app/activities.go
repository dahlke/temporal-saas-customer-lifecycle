package app

import (
	"context"
	"fmt"
)

func ChargeCustomer(ctx context.Context, customerID string, amount float64) (string, error) {
	result := fmt.Sprintf("Charged customer %s amount %.2f", customerID, amount)
	return result, nil
}

func UndoChargeCustomer(ctx context.Context, customerID string, amount float64) (string, error) {
	result := fmt.Sprintf("Charge for customer %s amount %.2f was undone", customerID, amount)
	return result, nil
}

func CreateAccount(ctx context.Context, email string, password string) (string, error) {
	result := fmt.Sprintf("Created account for %s", email)
	return result, nil
}

func UndoCreateAccount(ctx context.Context, email string) (string, error) {
	result := fmt.Sprintf("Account for %s was deleted", email)
	return result, nil
}

func CreateAdminUsers(ctx context.Context, emails []string) (string, error) {
	result := fmt.Sprintf("Created admin accounts for %v", emails)
	return result, nil
}

func UndoCreateAdminUsers(ctx context.Context, emails []string) (string, error) {
	result := fmt.Sprintf("Admin accounts for %v were deleted", emails)
	return result, nil
}

func SendClaimCodes(ctx context.Context, email string, code string) (string, error) {
	result := fmt.Sprintf("Sent claim code to user %s: %v", email, code)
	return result, nil
}

func SendWelcomeEmail(ctx context.Context, emails []string) (string, error) {
	result := fmt.Sprintf("Sent welcome email to %v", emails)
	return result, nil
}

func SendFeedbackEmail(ctx context.Context, emails []string) (string, error) {
	result := fmt.Sprintf("Sent feedback email to %v", emails)
	return result, nil
}
