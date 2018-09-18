import io

from csv2csoc import *


def test_process_input_print_output():
    fields = [0, 1, 3]
    input_csv = io.StringIO("D35,Šimerková Olga,SK Praga,PGP7551,C,16.31/,51,24/2, 67.55/1, 48.39/2,116.34/1\n")
    completed = process_data(fields, input_csv)
    assert completed == [['D35', 'Šimerková Olga', 'PGP7551', '1.00']]
    print_output(completed)
