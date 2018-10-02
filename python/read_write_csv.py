import csv
import json


def main():
    csvfile = open("./some_file.csv", newline='')
    reader = csv.DictReader(csvfile, delimiter='\t')
    match_records = {}
    errors = []
    for row in reader:
        if row['some_field'] not in match_records:
            match_records[row['some_field']] = {'some_field': row["another_field"]}
        else:
            match_records[row['some_field']]['id'].append(row['id'])

    if len(match_records) > 0:
        for sf_id, item in match_records.items():
            if not (item['some_item'] in item['another_item']):
                # do something or just write error
                errors.append(item)
                continue

        if len(errors) > 0:
            f = open("./errors.csv", newline='', mode='w')
            for val in errors:
                f.write(json.dumps(val))