#!/bin/sh

TEMPORAL_ENV="<cli_env_name>"

export TEMPORAL_ADDRESS=$(temporal env get --env ${TEMPORAL_ENV} --key address -o json | jq -r '.[].value')
export TEMPORAL_NAMESPACE=$(temporal env get --env ${TEMPORAL_ENV} --key namespace -o json | jq -r '.[].value')
export TEMPORAL_CERT_PATH=$(temporal env get --env ${TEMPORAL_ENV} --key tls-cert-path -o json | jq -r '.[].value')
export TEMPORAL_KEY_PATH=$(temporal env get --env ${TEMPORAL_ENV} --key tls-key-path -o json | jq -r '.[].value')
# These are the paths required for the `temporal` CLI
export TEMPORAL_TLS_CERT=$(temporal env get --env ${TEMPORAL_ENV} --key tls-cert-path -o json | jq -r '.[].value')
export TEMPORAL_TLS_KEY=$(temporal env get --env ${TEMPORAL_ENV} --key tls-key-path -o json | jq -r '.[].value')

# Optional
export TEMPORAL_TASK_QUEUE="LIFECYCLE_TASK_QUEUE"
export TEMPORAL_NEXUS_BILLING_TASK_QUEUE="SUBSCRIPTION_BILLING_TASK_QUEUE"
export TEMPORAL_NEXUS_BILLING_ENDPOINT="SUBSCRIPTION_BILLING_ENDPOINT"
export ENCRYPT_PAYLOADS=true
