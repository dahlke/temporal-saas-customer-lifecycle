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
