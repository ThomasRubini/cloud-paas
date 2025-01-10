#!/bin/sh
set -e

echo "Running custom init script to inject secrets into realm-config.json"

cd /tmp

jq "(.clients[] | select(.clientId == \"${OIDC_CLIENT_ID}\")).secret = \"${OIDC_CLIENT_SECRET}\"" base-realm-config.json > update.json
jq "(.clients[] | select(.clientId == \"${OIDC_CLIENT_ID}\")).secret = \"${OIDC_CLIENT_SECRET}\"" update.json > update2.json

IMPORT_DIR=/opt/keycloak/data/import/
mkdir -p $IMPORT_DIR
cp update2.json $IMPORT_DIR/realm-config.json

echo "Running custom script finished !"
export KEYCLOAK_IMPORT=$IMPORT_DIR/realm-config.json
/opt/keycloak/bin/kc.sh $@
