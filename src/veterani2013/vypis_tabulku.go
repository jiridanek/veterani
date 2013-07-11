package main

import (
	"fmt"
	"veterani2013/bodovani"
)

func main() {
	vypis()
}

func cara() {
	fmt.Println("------------------------------------------------------------------")

}
func vypis() {
	cara()
	fmt.Println("Přídělování bodů za umístění v nejvyšší V podkategorii ČPV")
	fmt.Printf("\n")
	table(0, 20, 22)
	fmt.Printf("\n")

	cara()
	fmt.Println("Přídělování bodů za umístění v druhé nejvyšší V podkategorii ČPV")
	fmt.Printf("\n")
	table(1, 10, 12)
	fmt.Printf("\n")

	cara()
	fmt.Println("Přídělování bodů za umístění v třetí nejvyšší V podkategorii ČPV")
	fmt.Printf("\n")
	table(2, 7, 7)
}

func table(kat, mu, mkz int) {
	fmt.Printf("        pořadí:")
	for u := 1; u <= mu; u++ {
		fmt.Printf("%4d", u)
	}
	fmt.Printf("\n")
	for kz := mkz; kz > 0; kz-- {
		if kz == mkz { // první řádek
			fmt.Printf("  kls.záv. ")
		} else {
			fmt.Printf("           ")
		}

		fmt.Printf("%-4d", kz)

		upto := kz
		if upto > mu {
			upto = mu
		}

		for u := 1; u <= upto; u++ {
			fmt.Printf("%4d", bodovani.Umisteni(kat, u, kz))
		}

		fmt.Printf("\n")
	}
}
