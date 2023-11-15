#!/usr/bin/env bash

set -e

OS=$(go env GOOS)
ARCH=$(go env GOARCH)
ROOT_DIR="$(cd $(dirname "${BASH_SOURCE}")/.. && pwd -P)"
PLUGIN="registry.terraform.io/pluralsh/plural"
PLUGIN_NAME="terraform-provider-$(basename "${PLUGIN}")"
PLUGIN_LOCATION="${ROOT_DIR}/build/${PLUGIN_NAME}"
VERSION=0.0.1
FILENAME="${PLUGIN_NAME}_v${VERSION}-${OS}-${ARCH}"
DEST_PATH="${HOME}/.terraform.d/plugins/${PLUGIN}/${VERSION}/${OS}_${ARCH}/${FILENAME}"

mkdir -p "$(dirname "${DEST_PATH}")"
cp "${PLUGIN_LOCATION}" "${DEST_PATH}"
echo "Installed ${PLUGIN} into ${DEST_PATH}"


