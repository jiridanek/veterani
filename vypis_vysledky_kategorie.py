import sqlalchemy
from collections import defaultdict

from sqlalchemy import create_engine
from typing import Dict, List
from dataclasses import dataclass


def category_from_regno(year, regno) -> int:
    last_two_year = year % 100
    first_two_regno = regno if regno < 100 else int(str(regno)[:2])

    if last_two_year < first_two_regno:
        last_two_year += 100

    return ((last_two_year - first_two_regno) // 5) * 5


def category_from_ID_age(kategorie: str, z_id: str) -> str:
    pohlavi = kategorie.split(",")[0][:1]
    regno = z_id.replace("|", "")
    cislo = category_from_regno(YEAR, int(regno[3:]))
    return pohlavi + str(cislo)


def get_race_results(conn, zid: str) -> Dict[int, int]:
    result = conn.execute("""
SELECT DISTINCT
  v.zavodnikID as v_zavodnikID,
  v.zavod as v_zavod,
  v.body as v_body
FROM vysledek v
WHERE zavodnikID=?""", zid)

    res = {}
    for row in result:
        res[row.v_zavod] = row.v_body
    return res


@dataclass
class ResultItem:
    kp_poradi: int
    p_poradi: int
    ap_poradi: int
    z_id: str
    z_prijmeni: str
    z_jmeno: str
    nzavodu: int
    kategorie: str
    s_body: int
    ap_scores: str

    @classmethod
    def fetch(cls, conn: sqlalchemy.engine.base.Connection):
        result: sqlalchemy.engine.result.ResultProxy = conn.execute("""
SELECT DISTINCT
  kp.poradi as kp_poradi,
  p.poradi as p_poradi,
  ap.poradi as ap_poradi,
  z.id as z_id,
  z.prijmeni as z_prijmeni,
  z.jmeno as z_jmeno,
  (SELECT COUNT(zavodnikID) FROM vysledek WHERE zavodnikID=z.id) as nzavodu,
  GROUP_CONCAT(kp.kat) as kategorie,
  s.body as s_body,
  ap.scores as ap_scores
FROM zavodnik z, soucet s, absporadi ap, poradi p, katporadi kp
WHERE z.id=s.zavodnikID AND z.id=ap.zavodnikID AND z.id=p.zavodnikID AND z.id=kp.zavodnikID
GROUP BY z.id
ORDER BY ap.poradi ASC
""")
        return [cls(**row) for row in result]


COLS: int = 35
HIGHLIGHT: int = 10
YEAR: int = 2018


def vypis_vysledky():
    engine = create_engine('sqlite:///db.sqlite', echo=True)
    conn = engine.connect()
    r = ResultItem.fetch(conn)

    results = defaultdict(list)
    for result in r:
        results[category_from_ID_age(result.kategorie, result.z_id)].append(result)

    with open(f"hodnoceni_cpv_{YEAR}_dle_kategorii_python.html", 'wt') as fp:
        print(vypis_html(results), file=fp)


def vypis_html(results: Dict[str, List[ResultItem]]):
    result = ["""
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
</head>
<body><pre>""",
              f"Pořadí Českého Poháru Veteránů {YEAR} dle kategorií\n\n"
              ]
    for k in sorted(results.keys()):  # over all categories
        result.append(f"     KATEGORIE {k}\n")
        result.append("poř.kat.  poř.ČPV abs.poř.  Reč         Jmeno             počet  body     body dle závodů\n")
        result.append("                                                         závodů celkem  ")
        for i in range(COLS):
            result.append(f"{i:2d}")
        result.append("\n")
        katporadi = 1
        for katporadi, l in enumerate(results[k], start=1):
            re_c = l.z_id.replace("|", "")
            prijmeni_jmeno = l.z_prijmeni + ', ' + l.z_jmeno
            result.append(
                f"{katporadi:7d} {l.p_poradi:6d} {l.ap_poradi:6d} {re_c:10s} {prijmeni_jmeno:<25s} {l.nzavodu:2d} {l.s_body:7d}    %s\n"
            )

    # katporadi, //db.Getkatporadi(l.Z_id, k), //l.Kp_poradi,
    # l.P_poradi,
    # l.Ap_poradi,
    # strings.Replace(l.Z_id, "|", "", -1),
    # l.Z_prijmeni+", "+l.Z_jmeno,
    # l.Nzavodu,
    # l.S_body,
    # sraces.String())

    result.append("\n")
    result.append("""</pre>
</body>
  </html>""")
    return ''.join(result)


if __name__ == '__main__':
    vypis_vysledky()
