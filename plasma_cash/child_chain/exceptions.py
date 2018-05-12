class InvalidTxSignatureException(Exception):
    """the signature of a tx is invalid"""


class InvalidBlockSignatureException(Exception):
    """the signature of a block is invalid"""


class PreviousTxNotFoundException(Exception):
    """previous transaction is not found"""


class TxAlreadySpentException(Exception):
    """the transaction is already spent"""


class TxAmountMismatchException(Exception):
    """tx input total amount is not equal to output total amount"""
