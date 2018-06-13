from flask import Flask
import json



app = Flask(__name__)

@app.route('/abci_info')
def block_number():
    resp = {
        "jsonrpc": "2.0",
        "id": "",
        "result": {
            "response": {
            "last_block_height": 4955,
            "last_block_app_hash": "p/FQg0Muf7IUMeQAE25T1JZDwFA="
            }
        }
    }
    return json.dumps(resp)

def main():
    app.debug = True
    app.run()

if __name__ == "__main__":
    main()