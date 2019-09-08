import json
import flask
import hashlib
import time


class Block:
    def __init__(self, index, appointments, timestamp, prev_hash, nonce=0):
        self.index = index
        self.appointments = appointments
        self.timestamp = timestamp
        self.prev_hash = prev_hash
        self.nonce = nonce

    def compute_hash(self):
        """ Just something simple for now """
        return hashlib.sha256(json.dumps(self.__dict__, sort_keys=True))


class BlockChain:
    def __init__(self):
        self.unconfirmed_apts = []
        self.chain = []

    def create_genesis_block(self):
        gen_block = Block(0, [], time.time(), "0")
        gen_block.hash = (
            gen_block.compute_hash()
        )  # todo : helper to make this in __init__
        self.chain.append(gen_block)

    def add_block(self, block, proof):
        pass

    def proof_of_work(self, block):
        pass

    def add_new_appointment(self, appointment):
        """ appointment is json format """
        pass

    def blocks_to_docs(self):
        """ Returns list of appointments in each block in chain to be stored in local db """
        pass

    def mine(self):
        pass


""" THIS BELOW SHOULD NOT BE GLOBAL SCOPE """

app = flask.Flask(__name__)


def consensus():
    """ if a longer chain exsists, gets the longest chain other miners, or nodes,
        and replaces local chain to longest.
    """
    pass


@app.route("/new_appointment", methods=["POST"])
def new_appointment():
    """ adds new appointment to chain """
    pass


@app.route("/mine", methods=["GET"])
def mine_unconfirmed():
    """ initiate mining on unconfirmed appointments """
    pass
