#!/bin/bash

set -e

sed -i -E 's/\$\.header\./\$\.headers\./g' "$PWD/client/pacts/pactflow-terraform-client-pactflow-application-saas.json"

echo "Done..."