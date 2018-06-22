from flask import Flask

from dependency_config import container


def create_app():
    app = Flask(__name__)

    # Create a child chain instance when creating a Flask app.
    container.get_child_chain()

    from . import server

    app.register_blueprint(server.bp)

    return app
