package input

import (
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
