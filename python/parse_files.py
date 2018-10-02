import os
import pymysql
import pymysql.cursors

gl = []


def main():
    root_dir = ''
    files = get_files(root_dir)
    tables = get_schema()
    for table_name in tables:
        parse_f(files, table_name)
        # r = list(sorted(set(gl)))
        # for nm in r:
        #     print(nm)


def get_files(root):
    d = []
    for path, subdirs, files in os.walk(root):
        for name in files:
            file = os.path.join(path, name)
            d.append(file)
    return d


def parse_f(files, phrase):
    for file in files:
        if file.endswith(".py") or file.endswith(".go"):
            with open(file, "r", encoding="utf-8") as f:
                searchlines = f.readlines()
            for i, line in enumerate(searchlines):
                if phrase in line:
                    # gl.append(file)
                    print("[%s] %s:%s :: %s" % (phrase, file, i, line))


def get_schema():
    connection = pymysql.connect(host='localhost',
                                 user='user',
                                 password='passwd',
                                 db='db_name',
                                 charset='utf8mb4',
                                 cursorclass=pymysql.cursors.SSCursor)
    tables = []
    with connection.cursor() as cursor:
        cursor.execute("SELECT table_name FROM information_schema.tables where table_schema = 'db_name';")

        for result in cursor.fetchall_unbuffered():
            tables.append(result[0])
    connection.close()
    return tables
