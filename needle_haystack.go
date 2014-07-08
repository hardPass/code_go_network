/*

This is a 'find-a-needle-in-haystack' thing.


*/

package main

import (
	"database/sql"
	"encoding/json"
	// "fmt"
	"bytes"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"time"
)

const (
	db_host     = "192.168.171.128:3306"
	db_user     = "root"
	db_password = "12345678"
)

var (
	db         *sql.DB
	haystack   []*Entity
	max_return = 12
)

type Spot struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Id       int    `json:"id"`
	Title    string `json:"title"`
}

func (this *Spot) String() string {
	b, _ := json.Marshal(this)
	return string(b)
}

func (this *Spot) StringPlain() string {
	return this.Province + " " + this.City + " " + this.Title
}

type Entity struct {
	Text  []byte
	Value interface{}
}

func NewEntity(spot *Spot) *Entity {

	ett := new(Entity)
	ett.Text = []byte(spot.StringPlain())
	ett.Value = spot
	return ett
}

func load() {
	log.Println("init...")

	var err error

	db, err = sql.Open("mysql", db_user+":"+db_password+"@tcp4("+db_host+")/res")
	if err != nil {
		log.Fatal(err)
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db connected:", db)

	// Prepare statement for reading data
	stmt, err := db.Prepare("select a.province,  a.city, a.id, a.title  from spot a order by a.province,a.city, a.title asc ")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	haystack = make([]*Entity, 0, 100000)
	for rows.Next() {
		spot := new(Spot)
		rows.Scan(&spot.Province, &spot.City, &spot.Id, &spot.Title)

		ett := NewEntity(spot)
		haystack = append(haystack, ett)
	}

	log.Printf("Loaded haystack %d .\n", len(haystack))
}

func search(needle []byte) (list []*Entity) {
	list = make([]*Entity, 0, 10)
	i := 0
	for _, v := range haystack {

		if bytes.Contains(v.Text, needle) {
			list = append(list, v)
			i++

			if i == max_return {
				return
			}
		}
	}

	return
}

type Tick struct {
	ns int64
}

func NewTick() *Tick {
	t := new(Tick)
	t.ns = time.Now().UnixNano()

	return t
}

func (this *Tick) Pick() int64 {
	tmp := this.ns
	this.ns = time.Now().UnixNano()
	return this.ns - tmp
}

func handler(w http.ResponseWriter, r *http.Request) {
	wd := r.FormValue("wd")
	log.Printf("----------wd param: %s.\n", wd)
	if len(wd) < 2 {
		return
	}

	t := NewTick()
	list := search([]byte(wd))
	log.Printf("search %s in %d ms!\n", wd, t.Pick()/(1000*1000))
	ret := make([]*Spot, 0, 10)
	for i, v := range list {
		spot, _ := v.Value.(*Spot)
		log.Printf("list[%d]: %v", i, spot.StringPlain())
		ret = append(ret, spot)
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(ret); nil != err {
		fmt.Fprintf(w, `{"error":"%s"}`, err)
	}
}

func main() {
	t := NewTick()
	load()
	log.Printf("loaded haystack in %d ms!\n", t.Pick()/(1000*1000))
	http.HandleFunc("/res/", handler)
	http.ListenAndServe(":8888", nil)

}
