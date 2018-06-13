# Development Dependencies

A patched version of web3.py is used because otherwise it does not work with Ganache due to issue #674. In addition, in order to be able to monitor events, PR #827, which is not merged yet. Pyethereum dependencies broke recently so we need to manually install a slightly older version of rlp encoding. Flask is used for server purposes.

On OSX + Homebrew
```
source /usr/local/bin/virtualenvwrapper.sh
```


```
mkvirtualenv erc721plasma --python=/usr/bin/python3.6
pip install -r requirements.txt
```

## Launch Plasma Chain

1. Make sure the contracts are deployed at the correct addresses (`npm run migrate:dev` in `server` directory)
2. Run `FLASK_APP=./child_chain FLASK_ENV=development flask run --port=8546` in one terminal. This will start a Plasma Chain instance which listens at `localhost:8546` and is also connected to the deployed contracts
3. Run `python demo.py`

TODO Should probably bundle these into makefiles, i.e. `make server` should launch the plasma chain.


## Testing

```
make test
```

## Linting

```
make lint
```
