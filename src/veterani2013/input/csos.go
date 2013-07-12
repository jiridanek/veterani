package input

import (
	"bufio"
	"fmt"
	//	"io"
	"log"
	"os"
	"strings"
	"veterani2013/types"
)

var _ = fmt.Printf

const CsocDisc string = "888.88"
const CsocAvg string = "999.99"

type Csos struct {
	Position    int
	Class       types.Class
	FamilyGiven string
	Regno       types.Regno
	License     string
	Result      string
}

func ReadCsos(fname string) []Csos {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	s := bufio.NewScanner(f)

	for s.Scan() {
		line := s.Text()
		if len(line) > 0 && line[0] == '-' {
			break
		}
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	rs := make([]Csos, 0)

	i := 1
	prevposition := 0
	prevresult := ""
	prevclass := ""
	for s.Scan() {
		line := []rune(s.Text())

		if len(line) != 49 {
			log.Printf("len(line) != 49, %d '%s'\n", len(line), string(line))
			continue
		}

		category := string(line[0:10])
		familyGiven := string(line[10:35])
		regno := string(line[35:42])
		license := string(line[42:43])
		result := string(line[43:49])

		// restart counter
		if prevclass != category {
			i = 1
			prevresult = ""
			prevposition = 0
		}

		position := i
		if prevresult == result {
			position = prevposition
		}

		r := Csos{Position: position,
			Class:       types.NewClass(strings.TrimSpace(category)),
			FamilyGiven: strings.TrimSpace(familyGiven),
			Regno:       types.NewRegno(strings.TrimSpace(regno)),
			License:     strings.TrimSpace(license),
			Result:      strings.TrimSpace(result)}

		rs = append(rs, r)

		if err := s.Err(); err != nil {
			log.Fatal(err)
		}
		i++
		prevclass = category
		prevresult = result
		prevposition = position
	}
	return rs
}
