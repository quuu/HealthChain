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

    @property
    def last_block(self):
        return self.chain[-1]

    @classmethod
    def is_valid_proof(cls, block, block_hash):
        return block_hash == block.compute_hash

    @classmethod
    def check_chain_validity(cls, chain):
        result = True
        previous_hash = 0  # genesus index
        for block in chain:
            block_hash = block.hash

            if (
                not cls.is_valid_proof(block, block_hash)
                or previous_hash != block.previous_hash
            ):
                result = False
                break

            # update prev hash for loop
            previous_hash = block_hash

        return result

    def create_genesis_block(self):
        # todo: db should have been initialized by this point, any other valid
        # blocks added should be added to db
        # how is seperate system going to coop with chain?
        gen_block = Block(0, [], time.time(), "0")
        gen_block.hash = (
            gen_block.compute_hash()
        )  # todo : helper to make this in __init__
        self.chain.append(gen_block)

    def add_block(self, block, proof):
        previous_hash = self.last_block.hash
        # two cases where block should not be added
        if previous_hash != block.previous_hash:
            return False
        if not BlockChain.is_valid_proof(block, proof):
            return False
        block.hash = proof
        self.chain.append(block)
        return True

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
