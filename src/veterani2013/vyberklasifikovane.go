package main

import (
  "fmt"
  "flag"
   "log"
  "veterani2013/input"
  "veterani2013/types"
  "veterani2013/output"
  "io"
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
for _,v := range rsklasif {
  klasmap[v.Regno] = true
}

fmt.Println("-")
rs := []input.Csos{}
for _,v := range rshlavni {
  _, found := klasmap[v.Regno]
  if found {
    rs = append(rs, v)
    l := []string{v.Class.String(),
      v.FamilyGiven,
      v.Regno.String(),
      v.License,
      v.Result}
    output.fprintfw(os.Stdout, l, []int{10,25,7,1,6})
  }
}



}