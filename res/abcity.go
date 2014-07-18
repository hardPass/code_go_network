/*

This is a 'find-a-needle-in-haystack' thing.


*/

package main

import (
	"encoding/json"
	"log"
)

func init() {
	registerLoadFunc("abcity", load_haystack_city_ab)
}

// data model for abroad city
type CityAb struct {
	Continent string
	Country   string
	City      string
	Id        int
}

func (this *CityAb) String() string {
	b, _ := json.Marshal(this)
	return string(b)
}

func (this *CityAb) Text() []byte {
	return []byte(this.Continent + " " + this.Country + " " + this.City)
}

func (this *CityAb) Json() interface{} {
	jo := struct {
		Id   int    `json:"id"`
		Text string `json:"text"`
	}{
		this.Id,
		this.Continent + " " + this.Country + " " + this.City,
	}

	return jo
}

func load_haystack_city_ab() []Entity {
	// Prepare statement for reading data
	stmt, err := db.Prepare("select id, country, city, continent from res_ab_city order by id asc")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatalf("--------------------------------", err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	haystack := make([]Entity, 0, 1000)
	for rows.Next() {
		mod := new(CityAb)
		rows.Scan(&mod.Id, &mod.Country, &mod.City, &mod.Continent)
		// log.Printf("%s \n", mod)
		haystack = append(haystack, mod)
	}

	return haystack
}
