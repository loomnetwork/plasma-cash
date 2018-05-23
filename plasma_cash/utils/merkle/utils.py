from web3.auto import w3

def keccak256(*args):
    hashes = list(args)
    bytes32 = ['bytes32'] * len(hashes)
    return w3.soliditySha3(bytes32, hashes)

def to_bytes(bits):
    # lord have mercy on my soul for this hack. 
    # There is probably a better way to store the slots that were 0/1 in a number
    # struct.pack('>Q', 0x200000000000000) <- hexstr does the job
    b = w3.toBytes(hexstr=hex(int(bits,2)))
    b  = b.rjust(8, b'\0') # pad to bytes needed
    return b

