package input

import (
	"bufio"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"veterani2013/types"
)

func IsValid(clubs map[string]bool, r types.Regno, c types.Class, p int, s string) bool {
	// club
	if !clubs[r.C] {
		return false
	}

	// class
	if c.B < 35 {
		return false
	}

	// position or status
	if p < 1 || s != "OK" {
		return false
	}

	return true
}

func GetValid(csos []Csos, oddily map[string]bool, vysledek Zavod) []Csos {
	rs := make([]Csos, 0)
	for _, r := range csos {
		// skip out of order
		if r.License == "M" {
			continue
		}
		// skip disc
		if r.Result == CsocDisc {
			continue
		}
		// skip foreigners, young ones, ...
		if !IsValid(oddily, r.Regno, r.Class, r.Position, "OK") {
			continue
		}
		// TODO? pokud jednorazovy zavod, preskoc vse krome kategorie B
		if vysledek.Jednorazovy {
			if r.Class.C != "B" {
				continue
			}
		}
		rs = append(rs, r)
	}
	return rs
}

func Nacti_oddily(fname string) map[string]bool {
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

func Nacti_zavod(dir, fname string) Zavod {
	parts := strings.Split(fname, "|")
	if len(parts) != 4 {
		log.Fatal("Format is number|attr|dd.mm.yyyy|name.suffix")
	}

	conv, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	cislo := int(conv)

	veteraniada := false
	jednorazovy := false
	switch parts[1] {
	case "v":
		veteraniada = true
	case "b":
		jednorazovy = true
	case "":
		// do nothing
	default:
		log.Fatalf("File %v: Attribute must be 'v', 'j' or '', was %v", fname, parts[1])
	}

	z := Zavod{cislo, veteraniada, jednorazovy, path.Join(dir, fname)}
	return z

}

type Zavod struct {
	Cislo       int
	Veteraniada bool
	Jednorazovy bool
	Fname       string
}
