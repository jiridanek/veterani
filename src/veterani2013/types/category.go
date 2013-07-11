package types

import (
  "fmt"
)

type Category struct {
  a string
  b int
  c string
}

func (c *Category) String() string {
  if b == 0 {
    return a
  }
  return fmt.Sprintf("%s%d%s", a,b,c)
}

func NewCategory(s string) *Category {
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
		return &Category{a, b, c}
	}

wrongformat:
	return &Category{s, 0, ""}
}