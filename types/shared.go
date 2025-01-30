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
	AccountName     string            `json:"account_name"`
	Price           float64           `json:"price"`
	Emails          []string          `json:"emails"`
	ClaimCodes      []ClaimCodeStatus `json:"claim_codes"`
	Progress        int               `json:"progress"`
	Status          string            `json:"status"`
	ChildWorkflowID string            `json:"child_workflow_id"`
	NexusWorkflowID string            `json:"nexus_workflow_id"`
	NexusNamespace  string            `json:"nexus_namespace"`
}

type LifecycleInput struct {
	AccountName string   `json:"account_name"`
	Emails      []string `json:"emails"`
	Price       float64  `json:"price"`
	Scenario    string   `json:"scenario"`
}
