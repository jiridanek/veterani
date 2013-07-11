package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"veterani2013/bodovani"
	"veterani2013/iof"
	"veterani2013/sql"
)

type Kategorie struct {
	a string
	b int
	c string
}

func NewKategorie(s string) Kategorie {
	//log.Println(s)

	if len(s) == 3 || len(s) == 4 {
		var a, c string
		var b int

		a = s[0:1]
		if a != "H" && a != "D" {
			goto divnyformat
		}

		conv, err := strconv.ParseInt(s[1:3], 10, 32)
		if err != nil {
			goto divnyformat
		}
		b = int(conv)
		if len(s) == 4 {
			c = s[3:4]
			switch c {
			case "A", "B", "C", "D", "E":
				// nedelej nic
			default:
				goto divnyformat
			}
		}
		return Kategorie{a, b, c}
	}

divnyformat:
	return Kategorie{s, 0, ""}
}

type Zavod struct {
	cislo       int
	attr bool
	fname       string
}

func nacti_zavod(dir, fname string) Zavod {
  parts := strings.Split(fname, "|")
		if len(parts) != 4 {
		  log.Fatal("Format is number|attr|dd.mm.yyyy|name.suffix")
		}

		conv, err := strconv.ParseInt(parts[0], 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		cislo := int(conv)

		attr := false
		
			switch(parts[1]) {
			  case "v":
			  attr = true
			  case "":
			    // do nothing
			  default:
				log.Fatal("Attribute must be ")
			}
			
		

		z := Zavod{cislo, attr, path.Join(dir, fname)}
		return z
		
}

func nacti_zavody(dir string) []Zavod {
	zavody := make([]Zavod, 0)
	fi, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, j := range fi {
		if j.IsDir() {
			continue
		}
		if path.Ext(j.Name()) != ".xml" {
		  continue
		}
		
z := nacti_zavod(dir, j.Name())
log.Printf("%#v\n", z)
		zavody = append(zavody, z)
	}
	return zavody
}

func nacti_oddily(fname string) map[string]bool {
	oddily := make(map[string]bool)
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		oddily[s.Text()] = true
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	return oddily
}

func main() {
	db := sql.NewDb("db.sqlite", true)
	//db := sql.NewDb(":memory:")
	naplndb(db, "clubs.txt", "2013/")
	zpracujdb(db)
	//celkove, kategorie := vysledky(db)
	db.Stop()
}

func zpracujdb(db sql.Db) {
	db.Vypocitejsoucty()
	db.Celkoveporadi()
	db.Poradi()
	db.Katporadi()
}

func naplndb(db sql.Db, foddily, dzavody string) {
	db.Createtables()
	// nacist oddily
	oddily := nacti_oddily(foddily)
	log.Printf("Oddilu: %d\n", len(oddily))
	log.Println("-------------")
	// nacist zavod po zavode
	vysledky := nacti_zavody(dzavody)
	log.Println("-------------")
	for _, vysledek := range vysledky {
		zavod := iof.Nacti_zavod(vysledek.fname)

		kategorie := make(map[Kategorie]bool)
		for _, r := range zavod.Results {
			kategorie[NewKategorie(r.Category)] = true
		}

		for k, _ := range kategorie {
			fmt.Printf("%v,", k)
		}
		fmt.Printf("\n")

		log.Printf("%#v\n", zavod.Event)
		log.Printf("Kategorii: %d\n", len(kategorie))
		log.Println("-------------")
		for _, r := range zavod.Results {
			log.Println(r.Category)

			kat := NewKategorie(r.Category)
			if kat.b < 35 {
				continue
			}
			katno := -1
			for _, k := range []string{"", "A", "B", "C", "D", "E"} {
				if kategorie[Kategorie{kat.a, kat.b, k}] {
					katno++
				}
				if k == kat.c {
					break
				}
			}
			if katno == -1 {
				log.Fatal("!!!BUG: katno!!!")
			}

			klaszav := len(r.PersonResults)

			log.Printf("Kategorie: %+v, Rank: %d, Zavodniku: %d", kat, katno, klaszav)
			log.Println("-------------")

			for _, p := range r.PersonResults {
				id := fmt.Sprintf("%s|%s", p.Person.Country, p.Person.Id)
				umisteni := p.Result.Position

				if umisteni < 1 || p.Result.Status.Value != "OK" {
					fmt.Printf("vynechavam umisteni: %v status: %v\n", umisteni, p.Result.Status.Value)
					continue
				}

				found := oddily[p.Person.Country]
				if !found {
					fmt.Printf("vynechavam Country: %v\n", p.Person.Country)
					continue
				}
				//printf

				b := bodovani.Ucast(katno) + bodovani.Umisteni(katno, umisteni, klaszav)
				if vysledek.attr {
					b += bodovani.Bonifikace(umisteni)
				}

				// 				if umisteni < 1 {
				// 				  fmt.Printf("umisteni: %v, body: %v\n", umisteni, b)
				// 				}

				db.Pridejvysledek(id, p.Person.Name.Given, p.Person.Name.Family, vysledek.cislo, b, kat.a, kat.b)
			}
		}
	}

	//m := make(map[string]bool)

	//fmt.Println(zavod.Event.Name)
	//fmt.Printf("%+v\n", zavod.Results)

	//fmt.Println(zavod.Results[0].PersonResults[0].Result.Status)
	//   fmt.Println(m)
}

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

func vypis_vysledky() {
	/* fmt.Println("Český pohár veteránů 2013")
	fmt.Println("")
	fmt.Println("D")
	fmt.Println("pořadí ČPV")
	*/
}
