package main

import (
	"flag"
	"fmt"
	"log"
	"veterani2013/input"
	"veterani2013/output"
	"veterani2013/types"
	//   "io"
	"os"
	//   "unicode/utf8"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 3 {
		log.Fatal("args: oddily klasif_csos hlavni_csos")

	}

	oddily := input.Nacti_oddily(args[0])

	klasif := input.Zavod{}
	hlavni := input.Zavod{}

	csosklasif := input.ReadCsos(args[1])
	rsklasif := input.GetValid(csosklasif, oddily, klasif)

	csoshlavni := input.ReadCsos(args[2])
	rshlavni := input.GetValid(csoshlavni, oddily, hlavni)

	klasmap := map[types.Regno]bool{}
	for _, v := range rsklasif {
		klasmap[v.Regno] = true
	}

	cs := make(map[types.Class][]input.Csos)
	for _, v := range rshlavni {
		_, found := klasmap[v.Regno]
		if found { // klasifikovany
			c := v.Class
			cs[c] = append(cs[c], v)
		}
	}

	keys := make([]types.Class, 0)
	for k, _ := range cs {
		keys = append(keys, k)
	}

	types.ClassBy(func(c1, c2 *types.Class) bool {
		if c1.A == c2.A {
			if c1.B == c2.B {
				return c1.C < c2.C
			} else {
				return c1.B < c2.B
			}
		} else {
			return c1.A < c2.A
		}
		return false
	}).Sort(keys)

	fmt.Println("-")

	// bug, pokud poslední v A má stejný čas
	// jako první v B, budou na děleném místě

	for _, c := range keys {
		for _, v := range cs[c] {
			l := []string{fmt.Sprintf("%s%d", v.Class.A, v.Class.B), // drop .C
				v.FamilyGiven,
				v.Regno.String(),
				v.License,
				v.Result}
			output.Fprintfw(os.Stdout, l, []int{10, 25, 7, 1, 6})
		}
	}
}
