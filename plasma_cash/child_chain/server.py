from flask import Blueprint, request
from dependency_config import container
import base64

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


@bp.route('/proof', methods=['GET'])
def get_proof():
    blknum = int(request.args.get('blknum'))
    uid = int(request.args.get('uid'))
    data = container.get_child_chain().get_proof(blknum, uid)
    return base64.b64encode(data) # proofs are binary so encoding b64 during transmission


@bp.route('/submit_block', methods=['POST'])
def submit_block():
    block = request.form['block']
    return container.get_child_chain().submit_block(block)


@bp.route('/send_tx', methods=['POST'])
def send_tx():
    tx = request.form['tx']
    return container.get_child_chain().send_transaction(tx)
