#!/bin/bash

set -e

sed -i -E 's/\$\.header\./\$\.headers\./g' "$PWD/client/pacts/terraform-client-delete-me-pactflow-application-saas.json"

echo "Done..."