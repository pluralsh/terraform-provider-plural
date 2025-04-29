#!/usr/bin/env bash

set -e

OS=$(go env GOOS)
ARCH=$(go env GOARCH)
ROOT_DIR="$(cd $(dirname "${BASH_SOURCE}")/.. && pwd -P)"
PLUGIN="registry.terraform.io/pluralsh/plural"
PLUGIN_NAME="terraform-provider-$(basename "${PLUGIN}")"
PLUGIN_LOCATION="${ROOT_DIR}/build/${PLUGIN_NAME}"
VERSION=$(curl -sL https://api.github.com/repos/pluralsh/terraform-provider-plural/tags | jq -r '.[0].name')
DESTINATION="${HOME}/.terraform.d/plugins/${PLUGIN}/${VERSION//v}/${OS}_${ARCH}/${PLUGIN_NAME}_${VERSION}-${OS}-${ARCH}"

mkdir -p "$(dirname "${DESTINATION}")"
mv "${PLUGIN_LOCATION}" "${DESTINATION}"
echo "Installed ${PLUGIN} ${VERSION} into ${DESTINATION}"
