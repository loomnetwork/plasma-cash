import requests

from .exceptions import RequestFailedException


class ChildChainService(object):
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
                'failed with response: {}'.format(response)
            )

    def get_current_block(self):
        end_point = '/block'
        response = self.request(end_point, 'GET')
        return response.text

    def get_block_number(self):
        end_point = '/blocknumber'
        response = self.request(end_point, 'GET')
        return int(response.text)

    def get_block(self, blknum):
        end_point = '/block/{}'.format(blknum)
        response = self.request(end_point, 'GET')
        return response.text

    def get_proof(self, blknum, slot):
        end_point = '/proof'
        params = {'blknum': blknum, 'slot': slot}
        response = self.request(end_point, 'GET', params=params)
        return response.text

    def get_tx_and_proof(self, blknum, slot):
        end_point = '/tx_proof'
        params = {'blknum': blknum, 'slot': slot}
        response = self.request(end_point, 'GET', params=params)
        return response.text

    def get_tx(self, blknum, slot):
        end_point = '/tx'
        params = {'blknum': blknum, 'slot': slot}
        response = self.request(end_point, 'GET', params=params)
        return response.text

    def submit_block(self):
        end_point = '/submit_block'
        response = self.request(end_point, 'POST')
        return int(response.text)

    def send_transaction(self, tx):
        end_point = '/send_tx'
        data = {'tx': tx}
        self.request(end_point, 'POST', data=data)
