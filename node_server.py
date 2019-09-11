import json
from flask import *
import hashlib
import time
from Crypto.Hash import SHA256


# initialize the flask app
app = Flask(__name__)



'''
Serves all static HTML files
'''
@app.route("/<string:page_name>/")
def hello(page_name):
    return render_template('%s.html' % page_name)


'''
Index page
'''
@app.route("/")
def root():
    return render_template('index.html')


'''
Route responsible for generating the hash-key
Should be removed for production

'''
@app.route("/hashkey")
def get_hash_key():
    first_name = request.args.get('first_name')
    if first_name is None:
        return "first_name not provided"
    last_name = request.args.get('last_name')
    if last_name is None:
        return "last_name not provided"
    country_code = request.args.get('country_code')
    if country_code is None:
        return "country_code not provided"
    unique_id = request.args.get('unique_id')
    if unique_id is None:
        return "unique_id not provided"

    key = SHA256.new(first_name + last_name + country_code + unique_id).hexdigest()
    return key


'''
Route to get the hashkey and compare with every
record in the database

'''
@app.route("/appointments")
def get_appointments():
    first_name = request.args.get('first_name')
    if first_name is None:
        return "first_name not provided"
    last_name = request.args.get('last_name')
    if last_name is None:
        return "last_name not provided"
    country_code = request.args.get('country_code')
    if country_code is None:
        return "country_code not provided"
    unique_id = request.args.get('unique_id')
    if unique_id is None:
        return "unique_id not provided"


    # perform for loop for every record and attempt to unhash
    return "some appointments"






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


