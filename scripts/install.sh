#!/bin/bash
#
# Usage:
#   $ curl -fsSL https://raw.githubusercontent.com/pactflow/terraform/master/scripts/install.sh | bash
# or
#   $ wget -q https://raw.githubusercontent.com/pactflow/terraform/master/scripts/install.sh -O- | bash
#

set -e
uname_output=$(uname)
case $uname_output in
  'Linux')
    linux_uname_output=$(uname -m)
    case $linux_uname_output in
      'x86_64')
        os='linux_amd64'
        ;;
      'i686')
        os='linux_386'
        ;;
      *)
        echo "Sorry, you'll need to install the terraform provider manually."
        exit 1
        ;;
    esac
    ;;
  'Darwin')
    os='darwin_amd64'
    ;;
  *)
  echo "Sorry, you'll need to install the terraform provider manually."
  exit 1
    ;;
esac

response=$(curl -s -v https://github.com/pactflow/terraform/releases/latest 2>&1)
tag=$(echo "$response" | grep -o "location: .*" | sed -e 's/[[:space:]]*$//' | grep -o "location: .*" | grep -o '[^/]*$')
mkdir -p ~/.terraform.d/plugins
curl -sLo ~/.terraform.d/plugins/terraform-provider-pact https://github.com/pactflow/terraform/releases/download/${tag}/terraform-provider-pact_${os}
