class InvalidTxSignatureException(Exception):
    """the signature of a tx is invalid"""


class CoinAlreadyIncludedException(Exception):
    """the coin uid has already been included in block"""


class InvalidBlockSignatureException(Exception):
    """the signature of a block is invalid"""


class PreviousTxNotFoundException(Exception):
    """previous transaction is not found"""


class TxAlreadySpentException(Exception):
    """the transaction is already spent"""


class TxAmountMismatchException(Exception):
    """tx input total amount is not equal to output total amount"""


class InvalidPrevBlockException(Exception):
    """prev block cannot be bigger than current block"""
