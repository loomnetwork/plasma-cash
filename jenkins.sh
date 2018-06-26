#!/bin/bash

mkdir -p gopath/bin ; true
export GOPATH=`pwd`/gopath

cd cash_test
make deps
make
make test
./plasmascash_tester ; true #remove true once finished


set -exo pipefail

# make test

cd server
npm install
npm run lint
npm run test

cd ../plasma_cash

#virtualenv --python=python3.5 .
#source bin/activate
export PATH="/var/lib/jenkins/.pyenv/bin:$PATH"
eval "$(pyenv init -)"
eval "$(pyenv virtualenv-init -)"
#pyenv virtualenv 3.6.0 general
pyenv global general

pip install -r requirements.txt
make lint
make test

cd ../
bash integration_test.sh
