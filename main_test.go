package main

import (
	"log"
	"testing"
	googlesheet "github.com/webserg/alertEngine/readGoogleSheet"
)


func TestGetSnameEq2(t *testing.T) {
	
	m, err := googlesheet.ReadTargetData()
	check(err)
	t.Logf("%s", "smty")
	log.Printf("%s", "sfdf")
	log.Println(m)
}

func TestGetPersone(t *testing.T) {
	// res, _ := http.Get("http://localhost:8080/people/2")
	// var person Person
	// _ = json.NewDecoder(res.Body).Decode(&person)
	// res.Body.Close()
	// if person.ID != "2"  {
	// t.Errorf("got pesone.ID = %s; want persone.ID=2", person.ID)
	// }
	// log.Printf("%s", person.ID)
}
