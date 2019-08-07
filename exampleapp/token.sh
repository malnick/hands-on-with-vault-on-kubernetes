#!/bin/bash -
set -Eeuo pipefail

JWT_TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
TOKEN=$(vault write -field token auth/kubernetes/login role=exampleapp-role jwt="${JWT_TOKEN}")

export VAULT_TOKEN=${TOKEN}
export VAULT_ADDR=http://localhost:8200
