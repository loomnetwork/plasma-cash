from flask import Blueprint, request
from dependency_config import container

bp = Blueprint('api', __name__)


@bp.route('/block', methods=['GET'])
def get_current_block():
    return container.get_child_chain().get_current_block()


@bp.route('/submit_block', methods=['POST'])
def submit_block():
    block = request.form['block']
    return container.get_child_chain().submit_block(block)


@bp.route('/send_tx', methods=['POST'])
def send_tx():
    tx = request.form['tx']
    return container.get_child_chain().send_transaction(tx)


@bp.route('/get_block', methods=['POST'])
def get_block():
    number = int(request.form['number'])
    return container.get_child_chain().get_block(number)
