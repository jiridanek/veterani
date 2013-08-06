package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io"
  "log"
  "os"
  "strconv"
  "strings"
  "unicode/utf8"
)

func main() {
  flag.Parse()
  args := flag.Args()
  if len(args) != 2 {
    log.Fatal("Args should be like:1,10,20 fname.csv")
  }
  
  cols := strings.Split(args[0], ",")
  log.Printf("Columns specified: %d\n", len(cols))
  coln := make([]int, 0)
  for _,v := range cols {
    i, err := strconv.ParseInt(v, 10, 32)
    if err != nil {
      log.Printf("Tried to convert '%s' to int:\n", v)
      log.Fatal(err)
    }
    coln = append(coln, int(i))
  }

  f, err := os.Open(args[1])
  if err != nil {
    log.Fatal(err)
  }
  
  r := csv.NewReader(f)
  r.FieldsPerRecord = len(coln)
fmt.Println("-")
  for {
    l, err := r.Read()
    if err == io.EOF {
      break
    }
    if err != nil {
      log.Fatal(err)
    }
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
      fmt.Printf("%s", v)
      //space
      for a := 0; a < coln[i]-lenv; a++ {
	fmt.Printf(" ")
      }
    }
    fmt.Printf("\n")
  }
}