package app

// TODO: merge this into the other shared files? Give it a different name?

const BillingServiceName = "billing-service"
const BillingOperationName = "charge-customer"

type BillingInput struct {
	AccountName string
	Price       float64
	Scenario    string
}

type BillingOutput struct {
	Message string
}
