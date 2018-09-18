import csv
import sys


def process_data(fields, input_csv):
    data = read_csv(input_csv)
    selected = select_fields(data, fields)
    completed = complete_time(selected)
    return completed


def read_csv(input_csv):
    r = csv.DictReader(input_csv, fieldnames=[], restkey=None)
    return [record[None] for record in r]


def select_fields(data, fields):
    return [_select_fields(datum, fields) for datum in data]


def _select_fields(datum, fields):
    try:
        return [datum[f] for f in fields]
    except IndexError:
        print("IndexError in", datum, file=sys.stderr)


def complete_time(data):
    next = 1
    result = []
    for datum in data:
        if len(datum) == 3:
            result.append(datum + [str(next) + ".00"])
            next += 1
        else:
            result.append(datum)
    return result


def print_output(data):
    print('--')
    for d in data:
        print(f'{d[0]:10s}{d[1]:25}{d[2]:8}{d[3]:>6}')


def main():
    fields = [int(f) for f in sys.argv[1].split(',')]
    with open(sys.argv[2], newline='') as input_csv:
        completed = process_data(fields, input_csv)
        print_output(completed)


if __name__ == '__main__':
    main()
