from tinydb import TinyDB, Query
import base64
from Crypto.Cipher import AES

import json

db = TinyDB('./db.json')


'''

returns encrypted appointment json using the key

'''
def encrypt_appointment(appointment, key):
    secret_key = key.rjust(32)

    cipher = AES.new(secret_key, AES.MODE_ECB)

    json_string = json.dumps(appointment)

    encoded = base64.b64encode(cipher.encrypt(json_string))



    return encoded

'''

returns True or False depending on if the message was decrypted successfully

'''
def decrypt_appointment(hashed, key):
    secret_key = key.rjust(32)

    cipher = AES.new(secret_key, AES.MODE_ECB)

    decoded = cipher.decrypt(base64.b64decode(hashed))

    try:
        obj = json.loads(decoded)
        return True
    except:
        return False


# test cases for the above functions


appointment = {
    'first': 'Andrew',
    'last': 'Qu',
    'date': 'today',
    'hospital': 'Albany Medical',
    'severity': 10
}

key = "jklxepewrwejnvsnkzcka"
print(encrypt_appointment(appointment,key))

wrong_key = "asdfavxcvqweradvasdfq"
print(decrypt_appointment(encrypt_appointment(appointment,key),key))

