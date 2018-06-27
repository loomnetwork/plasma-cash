#!/bin/bash

function cleanup {
    kill -9 $ganache_pid
    kill -9 $loom_pid
}

trap cleanup EXIT

mkdir loomtmp
cd loomtmp
# TODO download loom
#export LOOM_BIN=
$LOOM_BIN init

cd ../server
ganache_pid=$(npm run --silent migrate:dev)
echo 'Launched ganache' $ganache_pid

#cd ../loomtmp
#loom_pid=$($LOOM_BIN run)
#echo 'Launched loom' $loom_pid

# Wait for Ganache & Loom to spin up
#sleep 10

#cd ../loom_test
#./plasmas_tester
