package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"veterani2013/sql"
)

type vysledky struct {
	D []vysledek
	H []vysledek
}

type vysledky_po_kategoriich struct {
	kategorie []vysledek_po_kategoriich
}

type vysledek_po_kategoriich struct {
	Kategorie string
	vysledky  []vysledek
}

type vysledek struct {
	Pkat       int
	Pcpv       int
	Pabs       int
	Id         string
	Jmeno      string
	Nzav       int
	BodyCelkem int
	Body       []int
}

func main() {
	db := sql.NewDb("db.sqlite", false)
	vypis_vysledky(db)
	db.Stop()
}

func vypis_vysledky(db sql.Db) {
	r := db.Getresults()

	//fmt.Println(r[2])
	//return

	f, err := os.Create("hodnoceni_cpv_2013.txt")
	if err != nil {
		log.Fatal(err)
	}
	//defer f.Close()
	fmt.Fprintf(f, "\n")
	for _, k := range []string{"D", "H"} {
		fmt.Fprintf(f, "%s\n", k)
		fmt.Fprintf(f, "   pořadí  pořadí   ReČ             Jméno        počet    body\n")
		fmt.Fprintf(f, "    ČPV    absol.                               závodů\n")
		for _, l := range r {
			if l.Kategorie[0:1] != k {
				continue
			}
			fmt.Fprintf(f, "%7d %6d %10s %-24s %2d %7d    %s\n",
				l.P_poradi,
				l.Ap_poradi,
				strings.Replace(l.Z_id, "|", "", -1),
				l.Z_prijmeni+", "+l.Z_jmeno,
				l.Nzavodu,
				l.S_body,
				strings.Replace(l.Ap_scores, ",", " ", -1))
		}
		fmt.Fprintf(f, "\n")
	}
}
