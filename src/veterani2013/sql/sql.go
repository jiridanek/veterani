package sql

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

//func main() {}

type Db struct {
	*sql.Tx
	db *sql.DB
}

func NewDb(fname string, delete bool) Db {
	if delete {
		os.Remove(fname)
	}
	mdb, err := sql.Open("sqlite3", fname)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := mdb.Begin()
	if err != nil {
		log.Fatal(err)
	}

	return Db{tx, mdb}
}

func (db Db) Createtables() {
	sqls := []string{
		`CREATE TABLE zavodnik (
  id TEXT PRIMARY KEY,
  jmeno TEXT,
  prijmeni TEXT
);`,
		`CREATE TABLE vysledek (
  zavodnikID TEXT,
  zavod INTEGER,
  body INTEGER,
  kata TEXT,
  katb INTEGER
);`,
		`CREATE TABLE soucet (
  zavodnikID TEXT,
  body INTEGER
);`,
		`CREATE TABLE absporadi (
  zavodnikID TEXT,
  poradi INTEGER,
  scores TEXT
);`,
		`CREATE TABLE poradi (
  zavodnikID TEXT,
  poradi INTEGER,
  kata TEXT
);`,
		`CREATE TABLE katporadi (
  zavodnikID TEXT,
  poradi INTEGER,
  kat TEXT
);`,
		`CREATE VIEW vysledkovka_view AS
	  SELECT id, body, (SELECT GROUP_CONCAT( SUBSTR('00'||CAST(body AS TEXT),-2,2))
	    FROM (SELECT body
	      FROM vysledek
	      WHERE zavodnikID=id
	      ORDER BY body DESC)) AS scores
	  FROM zavodnik, soucet
	  WHERE id=zavodnikID
	  ORDER BY body DESC, scores DESC`,
	}
	for _, sql := range sqls {
		_, err := db.Exec(sql)
		if err != nil {
			log.Fatal("%q: %s\n", err, sql)
			return
		}
	}

}

func (db Db) Stop() {
	db.Commit()
	db.db.Close()
}

func (db Db) pridejzavodnika(id, jmeno, prijmeni string) {
	_, err := db.Exec("INSERT OR IGNORE INTO zavodnik (id, jmeno, prijmeni) values (?,?,?)", id, jmeno, prijmeni)
	if err != nil {
		log.Printf("%v, %v, %v", id, jmeno, prijmeni)
		log.Fatal(err)
		return
	}
}

func (db Db) Pridejvysledek(id, jmeno, prijmeni string, zavod, body int, kata string, katb int) {
	db.pridejzavodnika(id, jmeno, prijmeni)
	_, err := db.Exec("INSERT INTO vysledek (zavodnikID, zavod, body, kata, katb) VALUES (?, ?, ?, ?, ?)", id, zavod, body, kata, katb)
	if err != nil {
		//log.Printf("%v, %v, %v, %v", id, body, kata, katb)
		log.Fatal(err)
		return
	}
}

func (db Db) Vypocitejsoucty() {
	_, err := db.Exec("INSERT INTO soucet (zavodnikID, body) SELECT id,(SELECT SUM(body) FROM (SELECT body FROM vysledek WHERE zavodnikID=id ORDER BY body DESC LIMIT 10)) FROM zavodnik")
	if err != nil {
		log.Fatal(err)
		return
	}

}

func (db Db) Celkoveporadi() {
	//log.Println("celkove poradi")
	//rows, err := db.Query("SELECT id, body, (SELECT GROUP_CONCAT( SUBSTR('00'||CAST(body AS TEXT),-2,2)) FROM (SELECT body FROM vysledek WHERE zavodnikID=id ORDER BY body DESC)) AS scores FROM zavodnik, soucet WHERE id=zavodnikID ORDER BY body DESC, scores DESC")
	rows, err := db.Query("SELECT id, body, scores FROM vysledkovka_view")
	if err != nil {
		//log.Fatal("errrr")
		log.Fatal(err)
	}
	defer rows.Close()

	i := 1

	var prevporadi int
	var prevbody int
	var prevscores string

	for rows.Next() {
		var id string
		var body int
		var scores string
		rows.Scan(&id, &body, &scores)

		poradi := i

		if i != 1 {
			if body == prevbody && scores == prevscores {
				poradi = prevporadi
			}
		}

		//fmt.Println(id, name)
		_, err := db.Exec("INSERT INTO absporadi (zavodnikID, poradi, scores) VALUES (?, ?, ?)", id, poradi, scores)
		if err != nil {
			log.Fatal(err)
		}

		prevporadi = poradi
		prevbody = body
		prevscores = scores

		i++
	}
}

func (db Db) poradiveskupine() {
	type tempdata struct {
		i             int
		prevporadi    int
		prevabsporadi int
	}

	// DISTINCT: ve vysledek je kazdy zavodnik tolikrat, kolik bezel zavodu
	rows, err := db.Query("SELECT DISTINCT z.id, ap.poradi, v.kata FROM absporadi ap,vysledek v,zavodnik z WHERE ap.zavodnikID=id AND v.zavodnikID=id ORDER BY ap.poradi ASC")
	if err != nil {
		//log.Fatal("errrr")
		log.Fatal(err)
	}
	defer rows.Close()

	offsets := make(map[string]*tempdata)

	for rows.Next() {
		var id string
		var absporadi int
		var kata string
		rows.Scan(&id, &absporadi, &kata)

		var poradi int
		if offsets[kata] == nil {
			poradi = 1
			offsets[kata] = &tempdata{1, 0, 0}
		} else {
			poradi = offsets[kata].i
			if absporadi == offsets[kata].prevabsporadi {
				poradi = offsets[kata].prevporadi
			}
		}

		_, err := db.Exec("INSERT INTO poradi (zavodnikID, poradi, kata) VALUES (?, ?, ?)", id, poradi, kata)
		if err != nil {
			log.Fatal(err)
		}

		offsets[kata].prevporadi = poradi
		offsets[kata].prevabsporadi = absporadi
		offsets[kata].i++

	}
}

func (db Db) Poradi() {

	type tempdata struct {
		i             int
		prevporadi    int
		prevabsporadi int
	}

	// DISTINCT: ve vysledek je kazdy zavodnik tolikrat, kolik bezel zavodu
	rows, err := db.Query("SELECT DISTINCT z.id, ap.poradi, v.kata FROM absporadi ap,vysledek v,zavodnik z WHERE ap.zavodnikID=id AND v.zavodnikID=id ORDER BY ap.poradi ASC")
	if err != nil {
		//log.Fatal("errrr")
		log.Fatal(err)
	}
	defer rows.Close()

	offsets := make(map[string]*tempdata)

	for rows.Next() {
		var id string
		var absporadi int
		var kata string
		rows.Scan(&id, &absporadi, &kata)

		var poradi int
		if offsets[kata] == nil {
			poradi = 1
			offsets[kata] = &tempdata{1, 0, 0}
		} else {
			poradi = offsets[kata].i
			if absporadi == offsets[kata].prevabsporadi {
				poradi = offsets[kata].prevporadi
			}
		}

		_, err := db.Exec("INSERT INTO poradi (zavodnikID, poradi, kata) VALUES (?, ?, ?)", id, poradi, kata)
		if err != nil {
			log.Fatal(err)
		}

		offsets[kata].prevporadi = poradi
		offsets[kata].prevabsporadi = absporadi
		offsets[kata].i++

	}
}

func (db Db) Katporadi() {
	type tempdata struct {
		i             int
		prevporadi    int
		prevabsporadi int
	}

	// DISTINCT: ve vysledek je kazdy zavodnik tolikrat, kolik bezel zavodu
	rows, err := db.Query("SELECT DISTINCT z.id, ap.poradi, v.kata, v.katb FROM absporadi ap,vysledek v,zavodnik z WHERE ap.zavodnikID=id AND v.zavodnikID=id ORDER BY ap.poradi ASC")
	if err != nil {
		//log.Fatal("errrr")
		log.Fatal(err)
	}
	defer rows.Close()

	offsets := make(map[string]*tempdata)

	for rows.Next() {
		var id string
		var absporadi int
		var kata string
		var katb int
		rows.Scan(&id, &absporadi, &kata, &katb)
		kat := fmt.Sprintf("%s%d", kata, katb)

		var poradi int
		if offsets[kat] == nil {
			poradi = 1
			offsets[kat] = &tempdata{1, 0, 0}
		} else {
			poradi = offsets[kat].i
			if absporadi == offsets[kat].prevabsporadi {
				poradi = offsets[kat].prevporadi
			}
		}

		_, err := db.Exec("INSERT INTO katporadi (zavodnikID, poradi, kat) VALUES (?, ?, ?)", id, poradi, kat)
		if err != nil {
			log.Fatal(err)
		}

		offsets[kat].prevporadi = poradi
		offsets[kat].prevabsporadi = absporadi
		offsets[kat].i++

	}
}

type Result struct {
	Kp_poradi int
	P_poradi   int
	Ap_poradi  int
	Z_id       string
	Z_prijmeni string
	Z_jmeno    string
	Nzavodu    int
	Kategorie  string  // multiple comma separated categories in here
	S_body     int
	Ap_scores  string
}

func (db Db) Getresults() []Result {
	rows, err := db.Query(
		`SELECT DISTINCT
  kp.poradi,
  p.poradi,
  ap.poradi,
  z.id,
  z.prijmeni,
  z.jmeno,
  (SELECT COUNT(zavodnikID) FROM vysledek WHERE zavodnikID=z.id),
  GROUP_CONCAT(kp.kat),
  s.body,
  ap.scores
FROM zavodnik z, soucet s, absporadi ap, poradi p, katporadi kp
WHERE z.id=s.zavodnikID AND z.id=ap.zavodnikID AND z.id=p.zavodnikID AND z.id=kp.zavodnikID
GROUP BY z.id
ORDER BY ap.poradi ASC`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	results := make([]Result, 0)

	for rows.Next() {
	        var kp_poradi int
		var p_poradi int
		var ap_poradi int
		var z_id string
		var z_prijmeni string
		var z_jmeno string
		var nzavodu int
		var kategorie string
		var s_body int
		var ap_scores string

		//cols, _ := rows.Columns()
		//fmt.Println(cols)
		//rows.Scan(&kp_poradi)
		err = rows.Scan(&kp_poradi, &p_poradi, &ap_poradi, &z_id, &z_prijmeni, &z_jmeno, &nzavodu, &kategorie, &s_body, &ap_scores)
		if err != nil {
			log.Fatal(err)
		}

		r := Result{Kp_poradi: kp_poradi,
		  P_poradi: p_poradi,
			Ap_poradi:  ap_poradi,
			Z_id:       z_id,
			Z_prijmeni: z_prijmeni,
			Z_jmeno:    z_jmeno,
			Nzavodu:    nzavodu,
			Kategorie:  kategorie,
			S_body:     s_body,
			Ap_scores:  ap_scores}

		//fmt.Println(kp_poradi)

		results = append(results, r)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return results
}

func (db *Db) Getraceresults(id string) map[int]int {
  rows, err := db.Query(
		`SELECT DISTINCT
  v.zavodnikID,
  v.zavod,
  v.body
FROM vysledek v
WHERE zavodnikID=?`, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	res := make(map[int]int, 0)

	for rows.Next() {
	        var v_zavodnikID string
		var v_zavod int
		var v_body int

		//cols, _ := rows.Columns()
		//fmt.Println(cols)
		//rows.Scan(&kp_poradi)
		err = rows.Scan(&v_zavodnikID, &v_zavod, &v_body)
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println(kp_poradi)

		res[v_zavod] = v_body
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return res
}

func (db *Db) Getkatporadi(id, kat string) int {
  rows, err := db.Query(
		`SELECT kp.poradi
FROM katporadi kp
WHERE kp.zavodnikID=? AND kp.kat=?`, id, kat)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	rows.Next()
		var kp_poradi int

		err = rows.Scan(&kp_poradi)
		if err != nil {
			log.Fatal(err)
		}

		
	
	if err = rows.Err(); err != nil || rows.Next() {
		log.Fatal(err)
	}
		
	return kp_poradi
}

// kolik lidi bezelo ve vice kategoriich
// SELECT COUNT (id) as cnt FROM ( SELECT DISTINCT z.id as id, ap.poradi, v.katb FROM absporadi ap,vysledek v,zavodnik z WHERE ap.zavodnikID=id AND v.zavodnikID=id ) GROUP BY id ORDER BY cnt asc;
