#!/bin/bash
/usr/bin/token.sh
echo "token: $VAULT_TOKEN"
exampleapp "$@"
