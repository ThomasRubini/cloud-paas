# About keycloak config files

`base-realm-config.json` is the partial realm export available from the keycloak UI. Any further manual modifications from it must be done separately.
The way to apply the modifications is to use the custom script `start.sh` that will be run on each startup before keycloak.
An example of usage of this script is replacing a specific key inside the config (e.g. a client secret) to edit a key, or merging it with another json file to add new keys (e.g. add users).

This script uses jq to do that, which is why we use a custom keycloak image backed by the Dockerfile in this directory.
