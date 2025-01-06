#!/bin/sh

TEMPORAL_ENV="<cli_env_name>"

# These are the paths required for the `temporal` CLI
export TEMPORAL_ADDRESS=$(temporal env get --env ${TEMPORAL_ENV} --key address -o json | jq -r '.[].value')
export TEMPORAL_NAMESPACE=$(temporal env get --env ${TEMPORAL_ENV} --key namespace -o json | jq -r '.[].value')
export TEMPORAL_TLS_CERT=$(temporal env get --env ${TEMPORAL_ENV} --key tls-cert-path -o json | jq -r '.[].value')
export TEMPORAL_TLS_KEY=$(temporal env get --env ${TEMPORAL_ENV} --key tls-key-path -o json | jq -r '.[].value')
export TEMPORAL_API_KEY=$(temporal env get --env ${TEMPORAL_ENV} --key api-key -o json | jq -r '.[].value')

# Optional
export TEMPORAL_TASK_QUEUE="lifecycle-task-queue"
# TODO: make sure these are truly optional
export TEMPORAL_NEXUS_BILLING_TASK_QUEUE="SUBSCRIPTION_BILLING_TASK_QUEUE"
export TEMPORAL_NEXUS_BILLING_ENDPOINT="SUBSCRIPTION_BILLING_ENDPOINT"
export ENCRYPT_PAYLOADS=true
