package app

const OnboardingTaskQueue = "ONBOARDING_TASK_QUEUE"

type OnboardingWorkflowInput struct {
	AccountName string   `json:"account_name"`
	Emails      []string `json:"emails"`
}

type ClaimCodeStatus struct {
	Email     string
	Code      string
	IsClaimed bool
}

type OnboardingWorkflowState struct {
	AccountName string            `json:"account_name"`
	Emails      []string          `json:"emails"`
	ClaimCodes  []ClaimCodeStatus `json:"claim_codes"`
}

type AcceptClaimCodeInput struct {
	ClaimCode string `json:"claim_code"`
}
