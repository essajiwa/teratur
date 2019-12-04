#!/bin/bash
set -e

APP_DIR="$GOPATH/src/github.com/${GITHUB_REPOSITORY}"

mkdir -p ${APP_DIR}
cp -r ./ ${APP_DIR} && cd ${APP_DIR}

echo "== Install Dep =="
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
dep ensure -v -vendor-only

echo "== Install Static Check Tools =="
. /setup.sh

make test-all