# temporal-saas-customer-onboarding

## Onboarding Workflow

- ChargeCustomer
- CreateAccount
  - RefundCustomer if error
- CreateAdminUsers
  - RefundCustomer and DeleteAccount if error
- SendClaimCodes
  - Wait 2 minutes
  - ResendClaimCodesSignal to resend claim codes
  - AcceptClaimCodeUpdate to accept claim code
  - RefundCustomer, DeleteAccount, and DeleteAdminUsers if error
- SendWelcomeEmail
  - Wait 10 seconds to send the feedback email
- SendFeedbackEmail
  - Clear our Saga compensations
- ChargeCustomer on a loop every 10 seconds
  - CancelSubscriptionSignal to cancel

## Onboarding Signals

- CancelSubscriptionSignal
- ResendClaimCodesSignal

## Onboarding Updates

- AcceptClaimCodeUpdate

## Onboarding Queries

- GetState

## Setup

### Running and Configuring the Temporal Dev Server (Option #1)

If you are using the Temporal Dev Server, start the server with the `frontend.enableUpdateWorkflowExecution` config
option set to `true`, which will allow us to perform updates to our workflows.

```bash
temporal server start-dev --db-filename temporal.sqlite --dynamic-config-value frontend.enableUpdateWorkflowExecution=true
```

Before kicking off the starter or using the UI, make sure the custom search attributes have been
created. If you are using the Temporal dev server, use the `operator search-attribute create`
command.

```bash
temporal operator search-attribute create --namespace $TEMPORAL_NAMESPACE --name OnboardingStatus --type text
```

### Configuring Temporal Cloud (Option #2)

First, you will need to set the following environment variables if you are using Temporal Cloud.

```bash
export TEMPORAL_ADDRESS="<namespace>.<accountId>.tmprl.cloud:7233"
export TEMPORAL_CERT_PATH="/path/to/ca.pem"
export TEMPORAL_KEY_PATH="/path/to/ca.key"
export TEMPORAL_NAMESPACE="<namespace>"
export TEMPORAL_ONBOARDING_TASK_QUEUE="onboarding"
```

If you are using Temporal Cloud, the command will look a bit different, using `tcld namespace search-attributes-add`.
If you are not already logged into Temporal Cloud with `tcld` run `tcld login`.

```bash
tcld namespace search-attributes add -n $TEMPORAL_NAMESPACE --sa "OnboardingStatus=Text"
```

### Interacting with the Workflow

As a helper function, we can export the latest workflow id to a variable.

```bash
export LATEST_WORKFLOW_ID=$(temporal workflow list --limit 1  | awk 'NR==2 {print $2}')
```

The generic format to interact with the workflow is as follows:

```bash

temporal workflow signal \
    --workflow-id="<workflow-id>" \
    --name ResendClaimCodesSignal

temporal workflow signal \
    --workflow-id="<workflow-id>" \
    --name CancelSubscriptionSignal

temporal workflow update \
    --workflow-id="<workflow-id>" \
    --name AcceptClaimCodeUpdate \
    --input '{"claim_code": "<claim_code>"}'

temporal workflow query \
    --workflow-id="<workflow-id>" \
    --type="GetState"
```

And if you want to debug the most recent workflow, you can use the following commands:

```bash
temporal workflow signal \
    --workflow-id=$(temporal workflow list --limit 1  | awk 'NR==2 {print $2}') \
    --name ResendClaimCodesSignal \
    --input '{"email": "neil.dahlke@temporal.io"}'

temporal workflow signal \
    --workflow-id=$(temporal workflow list --limit 1  | awk 'NR==2 {print $2}') \
    --name CancelSubscriptionSignal \
    --input '{"email": "neil.dahlke@temporal.io"}'

temporal workflow update execute \
    --workflow-id=$(temporal workflow list --limit 1  | awk 'NR==2 {print $2}') \
    --name AcceptClaimCodeUpdate \
    --input '{"claim_code": "XXX"}'

temporal workflow query \
    --workflow-id=$(temporal workflow list --limit 1  | awk 'NR==2 {print $2}') \
    --type="GetState"
```

## TODO

- codec server / encrypt_payloads
- Update readme
