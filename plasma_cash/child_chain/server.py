from flask import Blueprint, jsonify, request

from dependency_config import container

bp = Blueprint('api', __name__)


@bp.route('/block/<blknum>', methods=['GET'])
def get_block(blknum):
    return container.get_child_chain().get_block(int(blknum))


@bp.route('/block', methods=['GET'])
def get_current_block():
    return container.get_child_chain().get_current_block()


@bp.route('/blocknumber', methods=['GET'])
def get_block_number():
    return str(container.get_child_chain().get_block_number())


@bp.route('/tx_proof', methods=['GET'])
def get_tx_and_proof():
    blknum = int(request.args.get('blknum'))
    slot = int(request.args.get('slot'))
    tx, proof = container.get_child_chain().get_tx_and_proof(
        int(blknum), int(slot)
    )
    return jsonify({'tx': tx, 'proof': proof})


@bp.route('/tx', methods=['GET'])
def get_tx():
    blknum = int(request.args.get('blknum'))
    slot = int(request.args.get('slot'))
    return container.get_child_chain().get_tx(blknum, slot)


@bp.route('/proof', methods=['GET'])
def get_proof():
    blknum = int(request.args.get('blknum'))
    slot = int(request.args.get('slot'))
    return container.get_child_chain().get_proof(blknum, slot)


@bp.route('/submit_block', methods=['POST'])
def submit_block():
    return container.get_child_chain().submit_block()


@bp.route('/send_tx', methods=['POST'])
def send_tx():
    tx = request.form['tx']
    return container.get_child_chain().send_transaction(tx)
