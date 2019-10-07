from tinydb import TinyDB, Query
import base64
from Crypto.Cipher import AES
import json

db = TinyDB("./db.json")


"""

returns encrypted appointment json using the key

"""


def encrypt_appointment(appointment, key):

    # makes key 32 bytes long
    secret_key = key.rjust(32)

    cipher = AES.new(secret_key, AES.MODE_ECB)

    # makes sure the dictionary is 1024 bytes long
    json_string = json.dumps(appointment).rjust(1024)

    encoded = base64.b64encode(cipher.encrypt(json_string))

    return encoded


"""

returns True or False depending on if the message was decrypted successfully

criteria for successful decrypt is valid json object
  - maybe change in future for more security

"""


def can_decrypt_appointment(hashed, key):

    # makes key 32 bytes long
    secret_key = key.rjust(32)

    cipher = AES.new(secret_key, AES.MODE_ECB)

    decoded = cipher.decrypt(base64.b64decode(hashed))

    # tries to make a json object, if failed, not decrypted
    try:
        obj = json.loads(decoded)
    except:
        return False

    return True


"""

Actually decrypts the appointment

SHOULD NOT BE CALLED BEFORE can_decrypt_appointment IS CALLED

returns the decrypted json object

"""


def decrypt_appointment(hashed, key):

    secret_key = key.rjust(32)

    cipher = AES.new(secret_key, AES.MODE_ECB)

    decoded = cipher.decrypt(base64.b64decode(hashed))

    obj = json.loads(decoded)

    return obj


"""

Goes through every record and sees if the key provided
unlocks that record

returns all the records that are relevant to that user

"""


def get_appointments(key):

    # to store decrypted appointments
    user_appointments = []

    # for everything in the database
    for item in db.all():

        if can_decrypt_appointment(item["_"], key):

            # store the decrypted
            user_appointments.append(decrypt_appointment(item["_"], key))

    return user_appointments


# test cases for the above functions

key = "jklxepewrwejnvsnkzcka"

appointment = {
    "first": "Andrew",
    "last": "Qu",
    "date": "today",
    "hospital": "Albany Medical",
    "severity": 6,
}

"""

print(len(encrypt_appointment(appointment,key).decode("utf-8")))
db.insert({'_':encrypt_appointment(appointment,key).decode("utf-8")})

print(decrypt_appointment(encrypt_appointment(appointment,key),key))
"""
wrong_key = "asdfavxcvqweradvasdfq"


print(get_appointments(wrong_key))
