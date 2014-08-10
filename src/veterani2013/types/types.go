package types

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"errors"
	"sort"
)

// type Count struct {
//   D map[interface{}]interface{}
//   Updatefn func 
// }
// 
// func (c *Count) Update() {
//   
//   i++
// }

type Class struct {
	A string
	B int
	C string
}

//sorting
type ClassBy func(p1, p2 *Class) bool
func (by ClassBy) Sort(classes []Class) {
	cs := &classSorter{
		classes: classes,
		by:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(cs)
}
type classSorter struct {
	classes []Class
	by      ClassBy
}
func (s *classSorter) Len() int {
	return len(s.classes)
}
func (s *classSorter) Swap(i, j int) {
	s.classes[i], s.classes[j] = s.classes[j], s.classes[i]
}
func (s *classSorter) Less(i, j int) bool {
	return s.by(&s.classes[i], &s.classes[j])
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

func ClassLess(c, d Class) bool {
  return c.A == d.A && c.B < d.B
}

type Regno struct {
	C string
	N string
	L string
}

func (s Regno) String() string {
  return fmt.Sprintf("%s%s%s", s.C, s.N, s.L)
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

//FIXME separate config package?
var current = 100 + 14 // for 2014

func (n Regno) ClassB() (int, error) {
  if len(n.N) != 4 {
    log.Println(n)
    return 0, errors.New("Wrong numerical part")
  }
  byear, err := strconv.ParseInt(n.N[:2], 10, 32)
  if err != nil {
    log.Fatal(err)
  }
  
  // 35, 40, 45, ..., 75, +++
  // FIXME this breaks for people born after 2000
  // < 35, does not matter
  c := ((current - int(byear)) / 5) * 5
  return c, nil
}