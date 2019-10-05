import random
import string, csv
import datetime

import names # pip install names

def gen_ssn(length=9):
    """ Makes random string of numbers for ssn """
    return ''.join(random.choice(string.digits) for i in range(len))

def add_country(row, country="US"):
    row["Country"] = country
    return row

def add_dob(row):
    row['dob'] = gen_dob(int(row["age"])).strftime("%Y-%m-%d")
    return row

def add_ssn(row, ssn):
    row['ssn'] = ssn
    return row

def add_full_name(row):
    row["full_name"] = names.get_full_name(gender=row["gender"])
    return row
    

def gen_dob(current_age):
    """ Returns a random date in datetime.date format with the corresponding 
        current age.
    """
    years = datetime.timedelta(days=current_age*365) 
    today_date = datetime.datetime.today()
    to_date = today_date - years 
    from_date = to_date - datetime.timedelta(years=1)
    int_delta = ((to_date - from_date) * 24 * 60 * 60 ) + ((to_date - from_date)).seconds
    random_seconds = random.randrange(random_seconds)
    return (start + datetime.timedelta(seconds=random_seconds)).datetime.date()

def update_row(row, ssn):


def get_starter_data(filepath):
    with open(filepath, 'r') as file:
        reader = csv.DictReader(file)
        return [e for e in reader]

    

def gen_dataset(starter_data):
    num_patients = len(starter_data)
    ssns = set()
    while len(ssns) != num_patients:
        ssns.add(gen_ssn)
    for 
    




