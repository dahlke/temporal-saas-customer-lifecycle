package types

type AcceptClaimCodeInput struct {
	ClaimCode string `json:"claim_code"`
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

type OnboardingWorkflowInput struct {
	AccountName string   `json:"account_name"`
	Emails      []string `json:"emails"`
}
