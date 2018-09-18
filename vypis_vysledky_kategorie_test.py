import sqlite3
import pytest
import unittest
import sqlalchemy.engine.result
from sqlalchemy import create_engine
from dataclasses import dataclass

from vypis_vysledky_kategorie import category_from_regno


def test_noorm():
    engine = create_engine('sqlite:///db.sqlite', echo=True)
    conn = engine.connect()
    r = ResultItem.fetch(conn)
    _ = r[0]





def test_output_format():
    with sqlite3.Connection("db.sqlite") as c:
        r = c.execute(query)
        b = r.fetchall()
        _ = b


def test_category_from_regno():
    year = 2018
    for a, b in [
        (0, 18541),  # dulezite jsou prvni dve cislice
        (0, 17),
        (0, 14),
        (5, 13),
        (15, 1),
        (15, 00),
        (15, 99),
        (30, 88),
    ]:
        assert a == category_from_regno(year, b)


def test_db_mame_vsechny_tabulky():
    with sqlite3.Connection("db.sqlite") as c:
        r = c.execute("SELECT name FROM sqlite_master WHERE type = 'table'").fetchall()
        assert sorted(r) == sorted((v,) for v in ['zavodnik', 'vysledek', 'soucet', 'absporadi', 'poradi', 'katporadi'])


def test_db_mame_vsechny_pohledy():
    with sqlite3.Connection("db.sqlite") as c:
        r = c.execute("SELECT name FROM sqlite_master where type = 'view'").fetchall()
        assert sorted(r) == sorted((v,) for v in ['vysledkovka_view'])

    for l in r:
        print(l)
    unittest.TestCase().assertCountEqual([1], [1, 2])
    r.close()
    c.close()
