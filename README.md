# temporal-saas-customer-lifecycle

This demo shows how to implement a naive customer lifecycle workflow using Temporal.

| Prerequisites      |    | __ | Features       |    | __ | Patterns            |    |
|:-------------------|----|----|----------------|----|----|---------------------|----|
| Network Connection | ✅ | __ | Timer          | ✅ | __ | Entity              | ✅ |
| Golang 1.18+       | ✅ | __ | Update         | ✅ | __ | Long-Running        | ✅ |
|                    |    | __ | Signal         | ✅ | __ | Fanout              | ✅ |
|                    |    | __ | Query          | ✅ | __ | Continue As New     |    |
|                    |    | __ | Data Converter | ✅ | __ | Manual Intervention | ✅ |
|                    |    | __ | mTLS Keys      | ✅ | __ | Saga                | ✅ |
|                    |    | __ | Retry          | ✅ | __ |                     |    |
|                    |    | __ | Data Converter | ✅ | __ |                     |    |
|                    |    | __ | Child Workflow | ✅ | __ |                     |    |
|                    |    | __ | Polyglot       | ✅ | __ |                     |    |
|                    |    | __ | Replay Tests   |    | __ |                     |    |
|                    |    | __ | API Keys       |    | __ |                     |    |

## Lifecycle Workflow

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
  - NOTE: this runs in a child workflow if the scenario is set to `SCENARIO_CHILD_WORKFLOW`
  - CancelSubscriptionSignal to cancel

## Lifecycle Signals

- CancelSubscriptionSignal
- ResendClaimCodesSignal

## Lifecycle Updates

- AcceptClaimCodeUpdate

## Lifecycle Queries

- GetState

## Lifecycle Scenarios

- Happy Path
- Flakey API
- Unexpected Bug
- Expected Error
- Child Workflow

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
temporal operator search-attribute create --namespace default --name LifecycleStatus --type text
```

### Configuring Temporal Cloud (Option #2)

First, you will need to set the following environment variables if you are using Temporal Cloud.

```bash
export TEMPORAL_ADDRESS="<namespace>.<accountId>.tmprl.cloud:7233"
export TEMPORAL_TLS_CERT="/path/to/ca.pem"
export TEMPORAL_TLS_KEY="/path/to/ca.key"
export TEMPORAL_NAMESPACE="<namespace>"
export TEMPORAL_LIFECYCLE_TASK_QUEUE="lifecycle"
```

If you are using Temporal Cloud, the command will look a bit different, using `tcld namespace search-attributes-add`.
If you are not already logged into Temporal Cloud with `tcld` run `tcld login`.

```bash
tcld namespace search-attributes add -n $TEMPORAL_NAMESPACE --sa "LifecycleStatus=Text"
```

### Using Encryption

This demo supports encrypting the data sent to Temporal Cloud. In order to do so, you'll need to set the following environment variable.

```bash
export ENCRYPT_PAYLOADS=true
```

### Running the Demo

With your environment variables set, you can run the worker.

```bash
go run worker/main.go
```

### Running the Starter

Also with your environment variables set, you can run the starter.

```bash
go run starter/main.go
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
    --input '{"email": "sa@temporal.io"}'

temporal workflow signal \
    --workflow-id=$(temporal workflow list --limit 1  | awk 'NR==2 {print $2}') \
    --name CancelSubscriptionSignal \
    --input '{"email": "sa@temporal.io"}'

temporal workflow update execute \
    --workflow-id=$(temporal workflow list --limit 1  | awk 'NR==2 {print $2}') \
    --name AcceptClaimCodeUpdate \
    --input '{"claim_code": "XXX"}'

temporal workflow query \
    --workflow-id=$(temporal workflow list --limit 1  | awk 'NR==2 {print $2}') \
    --type="GetState"
```

### Using the SA Shared Codec Server

In the Temporal UI, configure your Codec server to use `https://codec.tmprl-demo.cloud` and check
the "pass the user access token" box.
