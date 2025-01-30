# Nexus Temp README

This will get migrated to the official README when completed.

```bash
temporal operator namespace create --namespace neil-dahlke-dev
temporal operator namespace create --namespace neil-dahlke-dev-nexus-target

temporal operator nexus endpoint create \
  --name SUBSCRIPTION_BILLING_ENDPOINT \
  --target-namespace neil-dahlke-dev-nexus-target \
  --target-task-queue SUBSCRIPTION_BILLING_TASK_QUEUE
```

```bash
TEMPORAL_ENV="nexus-dev"

export NEXUS_BILLING_ADDRESS="neil-dahlke-dev-nexus-target.sdvdw.tmprl.cloud:7233"
export NEXUS_BILLING_NAMESPACE="neil-dahlke-dev-nexus-target.sdvdw"
export NEXUS_BILLING_TASK_QUEUE="subscription-billing-task-queue"
export NEXUS_BILLING_ENDPOINT="subscription-billing-endpoint"

temporal env set --env $TEMPORAL_ENV -k env -v $TEMPORAL_ENV
temporal env set --env $TEMPORAL_ENV -k address -v $NEXUS_BILLING_ADDRESS
temporal env set --env $TEMPORAL_ENV -k namespace -v $NEXUS_BILLING_NAMESPACE

export NEXUS_BILLING_ADDRESS="neil-dahlke-dev-nexus-target.sdvdw.tmprl.cloud:7233"
export NEXUS_BILLING_NAMESPACE="neil-dahlke-dev-nexus-target.sdvdw"
export NEXUS_BILLING_TASK_QUEUE="subscription-billing-task-queue"
export NEXUS_BILLING_ENDPOINT="subscription-billing-endpoint"
```
