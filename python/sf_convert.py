import string


def convert_salesforce_record(record):
    """
    Convert 15-digits record to 18-digits in SalesForce
    :param record: 15-digits record
    :return: 18-digits record 
    """
    if len(record) != 15:
        print("Incorrect salesforce record id: %s" % record)
        return None
    matches = list(string.ascii_uppercase)
    matches.extend(['0', '1', '2', '3', '4', '5'])
    suffix = ''
    for i in range(0, 3):
        f = 0
        for j in range(0, 5):
            if record[i*5+j] in list(string.ascii_uppercase):
                f += 1 << j
        suffix += matches[f]
    return record + suffix