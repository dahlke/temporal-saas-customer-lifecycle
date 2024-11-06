# temporal-saas-customer-onboarding

## Onboarding Steps

- ChargeCustomer
- CreateAccount
- CreateAdminUsers
- SendClaimCodes
- SendWelcomeEmail
  - Wait 30 seconds
- SendFeedbackEmail

```bash
temporal workflow signal \
    --workflow-id="<workflow-id>" \
    --name AcceptClaimCodeSignal

temporal workflow update \
    --workflow-id="<workflow-id>" \
    --name AcceptClaimCodeUpdate \
    --input '{"claim_code": "XXX"}'

temporal workflow query \
    --workflow-id="<workflow-id>" \
    --type="GetState"
```

temporal workflow query \
    --workflow-id="onboarding-workflow-5343895a-db9b-499e-91c4-ea5c0218ac2c" \
    --type="GetState"


temporal workflow update \
    --workflow-id="onboarding-workflow-5343895a-db9b-499e-91c4-ea5c0218ac2c" \
    --name AcceptClaimCodeUpdate \
    --input '{"claim_code": "XXX"}'