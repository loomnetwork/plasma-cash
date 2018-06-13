import requests
import json

from .exceptions import RequestFailedException


class LoomChainService(object):

    def __init__(self, base_url, verify=False, timeout=5):
        self.base_url = base_url
        self.verify = verify
        self.timeout = timeout

    def request(self, end_point, method, params=None, data=None, headers=None):
        url = self.base_url + end_point

        response = requests.request(
            method=method,
            url=url,
            params=params,
            data=data,
            headers=headers,
            verify=self.verify,
            timeout=self.timeout,
        )

        if response.ok:
            return response
        else:
            raise RequestFailedException(
                        'failed with response: {}'.format(response))

    def get_current_block(self):
        # Loom version 
        print("Loom current_block")
        end_point = '/block'
        response = self.request(end_point, 'GET')
        return response.text

    def get_block_number(self):
        end_point = '/abci_info'
        response = self.request(end_point, 'GET')
        j_resp = json.loads(response.text)
        print(response.text)
        return int(j_resp["result"]["response"]["last_block_height"])

    def get_block(self, blknum):
        end_point = '/block/{}'.format(blknum)
        response = self.request(end_point, 'GET')
        return response.text

    def get_proof(self, blknum, uid):
        end_point = '/proof'
        params = {'blknum': blknum, 'uid': uid}
        response = self.request(end_point, 'GET', params=params)
        return response.text

    def submit_block(self, block):
        end_point = '/submit_block'
        data = {'block': block}
        self.request(end_point, 'POST', data=data)

    def send_transaction(self, tx):
        end_point = '/send_tx'
        data = {'tx': tx}
        self.request(end_point, 'POST', data=data)
