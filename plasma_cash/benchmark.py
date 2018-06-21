from client.client import Client
from dependency_config import container

authority = Client(container.get_root('authority'),
                   container.get_token('authority'))
w3 = authority.root_chain.w3  # get w3 instance

extras = list(map(lambda index: Client(container.get_root('extras', index),
                                       container.get_token('extras', index)),
                  range(0, 900)))

for extra in extras:
    extra.register()
