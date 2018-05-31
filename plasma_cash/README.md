# Development Dependencies

A patched version of web3.py is used because otherwise it does not work with Ganache due to issue #674. In addition, in order to be able to monitor events, PR #827, which is not merged yet. Pyethereum dependencies broke recently so we need to manually install a slightly older version of rlp encoding. Flask is used for server purposes.

```
mkvirtualenv erc721plasma --python=/usr/bin/python3.6
pip install -r requirements.txt
```

## Run demo.py

This is not complete, I'm still figuring it out.

TOOD George can you please describe this process in detail?

```
{run ganache}
python main.py &
python demo.py
```

## Testing

```
make test
```
