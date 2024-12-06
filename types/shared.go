package types

type AcceptClaimCodeInput struct {
	ClaimCode string `json:"claim_code"`
}

type ClaimCodeStatus struct {
	Email     string
	Code      string
	IsClaimed bool
}

type LifecycleWorkflowState struct {
	AccountName string            `json:"account_name"`
	Price       float64           `json:"price"`
	Emails      []string          `json:"emails"`
	ClaimCodes  []ClaimCodeStatus `json:"claim_codes"`
}

type LifecycleWorkflowInput struct {
	AccountName string   `json:"account_name"`
	Emails      []string `json:"emails"`
	Price       float64  `json:"price"`
	Scenario    string   `json:"scenario"`
}
