package app

const OnboardingTaskQueue = "ONBOARDING_TASK_QUEUE"

type AcceptClaimCodeInput struct {
	ClaimCode string `json:"claim_code"`
}
