package main

import (
	"database/sql"
	"encoding/json"
	// "fmt"
	"bytes"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	db_host     = "192.168.171.128:3306"
	db_user     = "root"
	db_password = "12345678"
)

var (
	db *sql.DB
)

func init() {
	log.Println("init...")

	var err error

	db, err = sql.Open("mysql", db_user+":"+db_password+"@tcp4("+db_host+")/res")
	if err != nil {
		log.Fatal(err)
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db connected:", db)
}

type Spot struct {
	Province string
	City     string
	Id       int
	Title    string
}

func (this *Spot) String() string {
	b, _ := json.Marshal(this)
	return string(b)
}

func ToDisk(b []byte) {
	fname := "spots." + strconv.FormatInt(time.Now().Unix(), 10) + ".txt"
	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fullpath, _ := filepath.Abs(fname)
	log.Printf("Create file %s %s \n", fname, fullpath)
	f.Seek(0, 0)
	f.Write(b)
}

func main() {
	defer db.Close()

	// Prepare statement for reading data
	stmt, err := db.Prepare("select a.province,  a.city, a.id, a.title  from spot a order by a.province,a.city, a.title asc limit 1000, 20")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	provice_cities_map := make(map[string]map[string][]interface{})
	for rows.Next() {
		var spot Spot
		rows.Scan(&spot.Province, &spot.City, &spot.Id, &spot.Title)

		if spot.City == "" {
			spot.City = spot.Province
		}

		log.Printf("-->: %v \n", &spot)

		var city_spots_map = provice_cities_map[spot.Province]
		if city_spots_map == nil {
			city_spots_map = make(map[string][]interface{})
			provice_cities_map[spot.Province] = city_spots_map
		}

		if city_spots_map[spot.City] == nil {
			city_spots_map[spot.City] = make([]interface{}, 0, 20)
		}

		x := struct {
			Id    int    `json:"id"`
			Title string `json:"title"`
		}{
			spot.Id,
			spot.Title,
		}
		city_spots_map[spot.City] = append(city_spots_map[spot.City], &x)

	}

	b, err := json.Marshal(provice_cities_map)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	// log.Printf("-->: %v \n", string(b))

	ToDisk(b)

	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	ToDisk(out.Bytes())
	out.WriteTo(os.Stdout)

	log.Printf("done!")
}
