import pytest
# from client.client import Client
# from dependency_config import container
from client.loom_chain_service import LoomChainService
# from utils.utils import increaseTime
from flask import Flask, jsonify
import time
import subprocess




class TestLoomChainService():
  p = None

  def test_abci_info(self):
    clientxx = LoomChainService('http://localhost:5000')
    currentBlock = clientxx.get_block_number()
    print(currentBlock)
    assert currentBlock == 4955

  def setup_method(self, method):
    cmd = "python client_flask_test.py"
    self.p = subprocess.Popen("exec " + cmd, stdout=subprocess.PIPE, shell=True)
    time.sleep(1)

  def teardown_method(self, method):
    self.p.kill()
    
