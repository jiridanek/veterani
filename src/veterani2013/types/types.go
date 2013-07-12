package types

import (
	"fmt"
	"strconv"
	"strings"
)

type Class struct {
	A string
	B int
	C string
}

func (c Class) String() string {
	if c.B == 0 {
		return c.A
	}
	return fmt.Sprintf("%s%d%s", c.A, c.B, c.C)
}

func NewClass(s string) Class {
	if len(s) == 3 || len(s) == 4 {
		var a, c string
		var b int

		a = s[0:1]
		if a != "H" && a != "D" {
			goto wrongformat
		}

		conv, err := strconv.ParseInt(s[1:3], 10, 32)
		if err != nil {
			goto wrongformat
		}
		b = int(conv)
		if len(s) == 4 {
			c = s[3:4]
			switch c {
			case "A", "B", "C", "D", "E", "N":
				// do nothing
			default:
				goto wrongformat
			}
		}
		return Class{a, b, c}
	}

wrongformat:
	return Class{s, 0, ""}
}

type Regno struct {
	C string
	N string
	L string
}

func NewRegno(s string) Regno {
	c := s
	n := ""

	i := strings.IndexAny(s, "0123456789")
	if i == -1 && len(s) > 3 {
		i = 3
	}
	if i != -1 {
		c = s[0:i]
		n = s[i:]
	}
	return Regno{c, n, ""}
}
