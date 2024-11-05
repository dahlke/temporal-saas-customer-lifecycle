# temporal-saas-customer-onboarding

## Onboarding Steps

- ChargeCustomer
- CreateAccount
- CreateAdminUsers
- CreateSupportChannel
- SendClaimCodes
- SendWelcomeEmail
- WAIT
- SendFeedbackEmail

- Resend welcome email / claim codes Signal

temporal workflow signal \
    --workflow-id="<workflow-id>" \
    --name update_apply_decision

temporal workflow update \
    --workflow-id="<workflow-id>" \
    --name AcceptClaimCode \
    --input '{"claim_code": "XXX"}'
