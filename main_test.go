package main

import (
	"log"
	"testing"
	"strings"
	googlesheet "github.com/webserg/alertEngine/readGoogleSheet"
)


func TestGetSnameEq2(t *testing.T) {
	
	m, err := googlesheet.ReadTargetData()
	check(err)
	t.Logf("%s", "smty")
	log.Printf("%s", "sfdf")
	log.Println(m)
}

func TestGetTargetData(t *testing.T) {
	celebs := map[string]int{
		"Nicolas Cage":       50,
		"Selena Gomez":       21,
		"Jude Law":           41,
		"Scarlett Johansson": 29,
		"Scarlett Johansson1": 29,
		"Scarlett Johansson2": 29,
		"Scarlett Johansson3": 29,
	}

	
	keys := make([]string,0,len(celebs)) 
	for k := range celebs {
		keys= append(keys, k)
    }
	log.Println(strings.Join(keys, ","))
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
