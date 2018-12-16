#!/bin/bash

set -exo pipefail

REPO_ROOT=`pwd`

cd $REPO_ROOT/server
npm install
npm run lint
npm run test

cd $REPO_ROOT/loom_test
export TMP_GOPATH =/tmp/gopath-$BUILD_TAG
export GOPATH=/tmp/gopath-$BUILD_TAG:`pwd`
make clean
make deps
make demos
make contracts
make test

cd $REPO_ROOT
REPO_ROOT=`pwd` IS_JENKINS_ENV=true bash e2e_test.sh
