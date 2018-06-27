#!/bin/bash

set -exo pipefail

# make test

cd server
npm install
#npm run lint
#npm run test

cd ../plasma_cash

#virtualenv --python=python3.5 .
#source bin/activate
export PATH="/var/lib/jenkins/.pyenv/bin:$PATH"
eval "$(pyenv init -)"
eval "$(pyenv virtualenv-init -)"
#pyenv virtualenv 3.6.0 general
#pyenv global general

#pip install -r requirements.txt
#make lint
#make test

#cd ../
#bash integration_test.sh

# build the Go tester and run the unit tests
cd loom_test
export GOPATH=/tmp/gopath-$BUILD_TAG:`pwd`
make clean
make deps
make cli
make oracle
make test

#cd ../
#bash go_integration_test.sh