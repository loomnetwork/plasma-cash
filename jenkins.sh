#!/bin/bash

cd server
npm install
npm run lint
npm run test

cd ../plasma_cash
pip3 install -r requirements.txt

cd ..
bash integration_test.sh

# make lint
# make test
