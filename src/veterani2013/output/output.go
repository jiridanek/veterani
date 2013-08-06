package output

import (
  "fmt"
  "unicode/utf8"
  "io"
)

func Fprintfw(f io.Writer, l []string, coln []int) {
      // trim to the width
    for i,v := range l {
      lenv := utf8.RuneCountInString(v)
      if lenv > coln[i] {
	l[i] = string([]rune(v)[:coln[i]])
      }
    }
    // print out
    for i,v := range l {
      lenv := utf8.RuneCountInString(v)
      fmt.Fprintf(f, "%s", v)
      //space
      for a := 0; a < coln[i]-lenv; a++ {
	fmt.Fprintf(f, " ")
      }
    }
    fmt.Fprintf(f, "\n")
  }