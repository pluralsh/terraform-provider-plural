#!/usr/bin/env bash

set -e

OS=$(go env GOOS)
ARCH=$(go env GOARCH)
ROOT_DIR="$(cd $(dirname "${BASH_SOURCE}")/.. && pwd -P)"
PLUGIN="registry.terraform.io/pluralsh/plural"
PLUGIN_NAME="terraform-provider-$(basename "${PLUGIN}")"
PLUGIN_LOCATION="${ROOT_DIR}/build/${PLUGIN_NAME}"
VERSION=0.0.1
DESTINATION="${HOME}/.terraform.d/plugins/${PLUGIN}/${VERSION}/${OS}_${ARCH}/${PLUGIN_NAME}_v${VERSION}-${OS}-${ARCH}"

mkdir -p "$(dirname "${DESTINATION}")"
mv "${PLUGIN_LOCATION}" "${DESTINATION}"
echo "Installed ${PLUGIN} into ${DESTINATION}"


