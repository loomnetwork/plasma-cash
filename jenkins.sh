#!/bin/bash

set -euxo pipefail

cd server
npm install
npm run lint
npm run test

cd ../plasma_cash

virtualenv --python=python3.5 .
source bin/activate
pip install -r requirements.txt
make lint

cd ../
bash integration_test.sh

# make lint
# make test
