#!/bin/bash

function cleanup {
    kill -9 $flask_pid
}

trap cleanup EXIT

cd server
npm run migrate:dev

cd ../plasma_cash
FLASK_APP=./child_chain FLASK_ENV=development flask run --port=8546 &
flask_pid=$!

sleep 3

python demo.py
