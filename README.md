# temporal-saas-customer-onboarding

## Onboarding Steps

- ChargeCustomer
- CreateAccount
- CreateAdminUsers
- SendClaimCodes
- SendWelcomeEmail
- WAIT
- SendFeedbackEmail

```bash
temporal workflow signal \
    --workflow-id="<workflow-id>" \
    --name update_apply_decision

temporal workflow update \
    --workflow-id="<workflow-id>" \
    --name AcceptClaimCode \
    --input '{"claim_code": "XXX"}'
```
