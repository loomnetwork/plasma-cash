#!/bin/bash

function cleanup {
    kill -TERM $flask_pid
    kill -9 $ganache_pid
}

trap cleanup EXIT

cd server
ganache_pid=$(npm run --silent migrate:dev)
echo 'Launched ganache' $ganache_pid

cd ../plasma_cash
FLASK_APP=./child_chain FLASK_ENV=development flask run --port=8546 &
flask_pid=$!
echo 'Launched flask' $flask_pid

sleep 3

python benchmark.py
