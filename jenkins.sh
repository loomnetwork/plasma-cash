#!/bin/bash

cd server
npm install
npm run lint
npm run test

cd ../plasma_cash

virtualenv --python=python3.5 .
source bin/activate
pip install -r requirements.txt

cd ../
bash integration_test.sh

# make lint
# make test
