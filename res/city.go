/*

This is a 'find-a-needle-in-haystack' thing.


*/

package main

import (
	"encoding/json"
	"log"
)

func init() {
	registerLoadFunc("city", load_haystack_city)
}

// data model
type City struct {
	Province string
	City     string
	Id       int
}

func (this *City) String() string {
	b, _ := json.Marshal(this)
	return string(b)
}

func (this *City) Text() []byte {
	return []byte(this.Province + " " + this.City)
}

func (this *City) Json() interface{} {
	jo := struct {
		Id   int    `json:"id"`
		Text string `json:"text"`
	}{
		this.Id,
		this.Province + " " + this.City,
	}

	return jo
}

func load_haystack_city() []Entity {
	// Prepare statement for reading data
	stmt, err := db.Prepare("select province,city,id from res_city order by id asc ")
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
		mod := new(City)
		rows.Scan(&mod.Province, &mod.City, &mod.Id)
		// log.Printf("%s \n", mod)
		haystack = append(haystack, mod)
	}

	return haystack
}
