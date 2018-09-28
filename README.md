# Plasma Cash - ERC721/ERC20/ETH Supported

## Installation and Loom SDK integration

To see the integrations run against the Loom SDK you must download the loom SDK, refer to https://loomx.io/developers/docs/en/basic-install-osx.html. 

Under the `loom_test` directory there are all the samples in Go, that directly interact with Loom SDK. If you have cloned the repo and go dependencies are not found, in `loom_test` try:

```
export GOPATH=$GOPATH:`pwd`
```

## Loom integration  tests

```
cd server
npm install
npm run test
cd .. 

cd loom_test
make clean
make deps
make demos
make contracts
make test
cd ..

cd loom_js_test
yarn
yarn build
yarn copy-contracts
cd ..

LOOM_BIN=<ABSOLUTE_PATH_TO_LOOM> ./loom_integration_test.sh
```

## License info 

Please note different directories have different licenses. Please see license file in each folder respectively

* server directory -> GPLv2 [License](plasma_cash/License.md)
* loom_test directory  -> Loom Public License [License](loom_test/License.md)
