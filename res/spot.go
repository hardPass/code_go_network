/*

This is a 'find-a-needle-in-haystack' thing.


*/

package main

import (
	"encoding/json"
	"log"
)

func init() {
	registerLoadFunc("spot", load_haystack_spot)
}

// data model
type Spot struct {
	Province string
	City     string
	Title    string
	CityId   int
	Id       int
}

func (this *Spot) String() string {
	b, _ := json.Marshal(this)
	return string(b)
}

func (this *Spot) Text() []byte {
	return []byte(this.Province + " " + this.City + " " + this.Title)
}

func (this *Spot) Json() interface{} {
	jo := struct {
		CityId int    `json:"city_id"`
		Id     int    `json:"id"`
		Text   string `json:"text"`
	}{
		this.CityId,
		this.Id,
		this.Province + " " + this.City + " " + this.Title,
	}

	return jo
}

func load_haystack_spot() []Entity {
	// Prepare statement for reading data
	stmt, err := db.Prepare("select a.province,  a.city, a.city_id,a.id, a.title  from res_spot a order by a.province,a.city, a.title asc ")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatalf("--------------------------------", err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	haystack := make([]Entity, 0, 100000)
	for rows.Next() {
		mod := new(Spot)
		rows.Scan(&mod.Province, &mod.City, &mod.CityId, &mod.Id, &mod.Title)
		// log.Printf("%s \n", mod)
		haystack = append(haystack, mod)
	}

	return haystack
}
