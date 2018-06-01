#!/usr/bin/env bash

function cleanup {
    kill -9 $ganache_pid
}

trap cleanup exit

testrpc_port=8545

testrpc_running() {
  nc -z localhost "$testrpc_port"
}

start_testrpc() {
  ganache-cli -a 15 -i 15 --blocktime 15 --gasLimit 50000000 -e 10000000000000000000000 -m gravity top burden flip student usage spell purchase hundred improve check genre > /dev/null 2>&1 &
  ganache_pid=$!
  echo "ganache-cli started with pid $ganache_pid"
}

if testrpc_running; then
  echo "Using existing testrpc instance at port $testrpc_port"
else
  echo "Starting our own testrpc instance at port $testrpc_port"
  start_testrpc
fi
