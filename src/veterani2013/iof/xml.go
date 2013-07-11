package iof

import (
	"encoding/xml"
      //"fmt"
	"log"
	"os"
  
)

//var _ = fmt.Printf

func Nacti_zavod(fname string) TResultList {
  f, err := os.Open(fname)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  d := xml.NewDecoder(f)
  d.Strict = true;
  
  var zavod TResultList 
  err = d.Decode(&zavod)
  
  if err != nil {
    log.Fatal(err)
  }
  return zavod
}

type TResultList struct {
  Event TEvent `xml:"Event"`
  Results []TClassResult `xml:"ClassResult"`
}

type TEvent struct {
  Id string `xml:"EventId"`
  Name string `xml:"Name"`
  // ...
}

type TClassResult struct {
  Category string `xml:"ClassShortName"`
  PersonResults []TPersonResult `xml:"PersonResult"`
}

type TPersonResult struct {
  Person TPerson `xml:"Person"`
  Result TResult `xml:"Result"`
}

type TPerson struct {
  Name TPersonName `xml:"PersonName"`
  Id string `xml:"PersonId"`
  Country string `xml:"CountryId"`
}

type TPersonName struct {
  Given string `xml:"Given"`
  Family string `xml:"Family"`
}

type TResult struct {
  StartTime string `xml:"StartTime"`
  FinishTime string `xml:"FinishTime"`
  Time string `xml:"Time"`
  Position int `xml:"ResultPosition"`
  Status TCompetitorStatus `xml:"CompetitorStatus"`
  // ...
}

type TCompetitorStatus struct {
  Value string `xml:"value,attr"`
}