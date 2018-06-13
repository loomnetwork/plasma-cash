#!/bin/bash

cd server
npm install
npm run lint
npm run test

cd ../plasma_cash

sudo apt install python3.6
sudo pip install virtualenv
virtualenv --python=python3.6 .
source bin/activate
pip install -r requirements.txt

cd ../
bash integration_test.sh

# make lint
# make test
