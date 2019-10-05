import random
import string
import csv
import datetime
import sys

import names # pip install names

def gen_ssn(length=9):
    """ Makes random string of numbers for ssn """
    return ''.join([random.choice(string.digits) for i in range(length)])

def add_country(row, country="US"):
    row["Country"] = country
    # return row

def add_dob(row):
    row['dob'] = gen_dob(int(float(row["age"]))).strftime("%Y-%m-%d")
    # return row

def add_ssn(row, ssn):
    row['ssn'] = ssn
    # return row

def add_full_name(row):
    row["full_name"] = names.get_full_name(gender=row["gender"])
    # return row
    

def gen_dob(current_age):
    """ Returns a random date in datetime.date format with the corresponding 
        current age.
    """
    years = datetime.timedelta(days=current_age*365) 
    today_date = datetime.datetime.today()
    to_date = today_date - years 
    from_date = to_date - datetime.timedelta(days=365)
    int_delta = ((to_date - from_date).days * 24 * 60 * 60 ) + (to_date - from_date).seconds
    random_seconds = random.randrange(int_delta)
    return (from_date + datetime.timedelta(seconds=random_seconds)).date()

def update_row(row, ssn):
    add_country(row)
    add_dob(row)
    add_full_name(row)
    add_ssn(row, ssn)
    # return row


def get_starter_data(filepath):
    with open(filepath, 'r') as file:
        reader = csv.DictReader(file)
        return [e for e in reader]

    

def gen_dataset(starter_data):
    num_patients = len(starter_data)
    ssns = set()
    while len(ssns) != num_patients:
        ssns.add(gen_ssn())
    for row in starter_data:
        new_ssn = random.sample(ssns, 1)[0]
        ssns.remove(new_ssn)
        # row = update_row(row, new_ssn)
        update_row(row, new_ssn)
    return starter_data

def main():
    data_fp = sys.argv[1]
    starter_data = get_starter_data(data_fp)
    hc_dataset = gen_dataset(starter_data)
    with open("hc_data.csv", 'w') as file:
        fields = list(hc_dataset[0].keys())
        writer = csv.DictWriter(file, fieldnames=fields)
        writer.writeheader()
        writer.writerows(hc_dataset)

    
if __name__ == "__main__":
    main()





