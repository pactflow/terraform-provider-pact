#!/bin/bash

set -e
cat $PWD/client/pacts/terraform-client-pactflow-application-saas.json | grep $.header
sed -i.bak -E "s/\$\.header\./\$\.headers\./g" "$PWD/client/pacts/terraform-client-pactflow-application-saas.json"
cat $PWD/client/pacts/terraform-client-pactflow-application-saas.json | grep $.header

echo "Done..."