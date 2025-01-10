#!/bin/sh
set -e

echo "Running custom init script to inject secrets into realm-config.json"

cd /tmp

# We can't put this change in a separate config file and merge it because it's using and array and idk how to merge arrays based on a key (clientId)
jq "(.clients[] | select(.clientId == \"${OIDC_CLIENT_ID}\")).secret = \"${OIDC_CLIENT_SECRET}\"" base-realm-config.json > tmp1.json
# This change on the other hand can, since the "users" key doesn't even exist in the base config
jq -s '.[0] * .[1]' tmp1.json users-realm-config.json > tmp2.json

IMPORT_DIR=/opt/keycloak/data/import/
mkdir -p $IMPORT_DIR
cp tmp2.json $IMPORT_DIR/realm-config.json

echo "Running custom script finished !"
export KEYCLOAK_IMPORT=$IMPORT_DIR/realm-config.json
/opt/keycloak/bin/kc.sh $@
