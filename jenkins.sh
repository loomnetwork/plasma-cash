#!/bin/bash

cd server
npm install
npm run lint
npm run test

# cd ../plasma_cash
# make lint
# make test
