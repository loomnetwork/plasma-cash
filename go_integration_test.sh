#!/bin/bash

function cleanup {
    kill -9 $ganache_pid
    kill -9 $loom_pid
}

REPO_ROOT=`pwd`
LOOM_DIR=/tmp/loom-plasma-$BUILD_TAG
BUILD_NUMBER=196

mkdir -p $LOOM_DIR
cd $LOOM_DIR
wget https://private.delegatecall.com/loom/linux/build-$BUILD_NUMBER/loom
chmod +x loom
export LOOM_BIN=`pwd`/loom
cp $REPO_ROOT/loom_test/loom-test.yml $LOOM_DIR/loom.yml
$LOOM_BIN init
echo 'Loom DAppChain initialized in ' $LOOM_DIR

trap cleanup EXIT

cd $REPO_ROOT/server
ganache_pid=$(npm run --silent migrate:dev)
echo 'Launched ganache' $ganache_pid

cd $LOOM_DIR
loom_pid=$($LOOM_BIN run)
echo 'Launched loom' $loom_pid

# Wait for Ganache & Loom to spin up
sleep 10

cd $REPO_ROOT/loom_test
./plasmascash_tester

kill -9 $loom_pid
rm -rf $LOOM_DIR
