package bodovani

import (
	"fmt"
	"log"
	"math"
	"veterani2013/types"
)

var _ = fmt.Printf

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// func round(a float64) int {
// 	whole := int(a)
// 	fractional := a - float64(whole)
// 	if math.Abs(fractional) >= 0.5 {
// 		if whole < 0 {
// 			return whole - 1
// 		}
// 		return whole + 1
// 	}
// 	return whole
// }

func Score(crank, pos, ncontestants int, veteraniada bool) int {
	s := Ucast(crank) + Umisteni(crank, pos, ncontestants)
	if veteraniada {
		s += Bonifikace(pos)
	}

	return s

}

func Ucast(kat int) int {
	switch kat {
	case 0:
		return 5
	case 1:
		return 3
	case 2:
		return 1
	default :
	  return 0
	}

}

// Za účast 2 body.
// Za umístění - vítěz 5 bodů a každý další v pořadí o bod méně.
func Bonifikace(u int) int {
	return 2 + max(5-(u-1), 0)
}

func Umisteni(kat, u, kz int) int {
	switch kat {
	case 0:
		var vitez float64
		var body float64
		vitez = 15.0
		body = vitez - float64(u-1)*(15.0/math.Min(0.75*float64(kz), 15.0))
		if kz < 8 {
			body *= float64(kz) / 8.0
		}
		// https://code.google.com/p/go/issues/detail?id=4594
		return max(int(body+0.5), 0)

	case 1:
		var vitez int
		var body int
		vitez = 6
		if kz < 10 {
			vitez = max(kz-4, 1)
		}

		body = vitez - (u - 1)
		return max(body, 0)

	case 2:
		var vitez int
		var body int
		vitez = 3
		if kz < 7 {
			vitez = kz - 4
		}

		body = vitez - (u - 1)
		return max(body, 0)
	default:
	  return 0
	}
}

func SubClassRank(cs map[types.Class]bool, c types.Class) int {
	r := -1
	for _, j := range []string{"A", "B", "C", "D", "E", ""} {
		if cs[types.Class{c.A, c.B, j}] {
			r++
		}
		if j == c.C {
			break
		}
	}
	// FIXME: move this check elsewhere
	if r == -1 {
		log.Fatal("!!!BUG: katno!!!")
	}
	return r

}
