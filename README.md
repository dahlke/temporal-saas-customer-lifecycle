# temporal-saas-customer-onboarding

## Onboarding Steps

- ChargeCustomer
- CreateAccount
- CreateAdminUsers
- SendClaimCodes
- SendWelcomeEmail
  - Wait 10 seconds
- SendFeedbackEmail

```bash
export LATEST_WORKFLOW_ID=$(temporal workflow list --limit 1  | awk 'NR==2 {print $2}')

temporal workflow signal \
    --workflow-id="<workflow-id>" \
    --name ResendClaimCodesSignal

temporal workflow update \
    --workflow-id="<workflow-id>" \
    --name AcceptClaimCodeUpdate \
    --input '{"claim_code": "XXX"}'

temporal workflow query \
    --workflow-id="<workflow-id>" \
    --type="GetState"
```

## Debugging

```bash
temporal workflow signal \
    --workflow-id=$(temporal workflow list --limit 1  | awk 'NR==2 {print $2}') \
    --name ResendClaimCodesSignal

temporal workflow update \
    --workflow-id=$(temporal workflow list --limit 1  | awk 'NR==2 {print $2}') \
    --name AcceptClaimCodeUpdate \
    --input '{"claim_code": "XXX"}'

temporal workflow query \
    --workflow-id=$(temporal workflow list --limit 1  | awk 'NR==2 {print $2}') \
    --type="GetState"
```
