/*

This is a 'find-a-needle-in-haystack' thing.


*/

package main

import (
	"database/sql"
	"encoding/json"
	// "fmt"
	"bytes"

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
	haystacks  = make(map[string][]Entity, 10)
	max_return = 20

	loadFuncMap = make(map[string]loadFunc, 0)
)

type loadFunc func() []Entity

func registerLoadFunc(name string, f loadFunc) {
	loadFuncMap[name] = f
	log.Printf("register loading func for : %s \n", name)
}

func loadHaystacks() {
	t := NewTick()
	for k, v := range loadFuncMap {

		list := v()
		haystacks[k] = list
		log.Printf("loaded a haystack for %s [len=%d] in %d ms!\n", k, len(list), t.Pick()/(1000*1000))
	}

}

type Entity interface {
	Json() interface{} // for output
	Text() []byte      // for search
}

func init_db() {
	log.Println("init_db...")

	var err error

	db, err = sql.Open("mysql", db_user+":"+db_password+"@tcp4("+db_host+")/hlg")
	if err != nil {
		log.Fatal(err)
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	log.Println("db connected:", db)
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

func printParams(r *http.Request) {

	r.ParseForm()
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			log.Printf("url params ------------ %s = %s\n", k, v[0])
		}
	}
}

func search(haystack []Entity, needle []byte) (list []Entity) {
	list = make([]Entity, 0, 10)
	i := 0
	for _, v := range haystack {

		if bytes.Contains(v.Text(), needle) {
			list = append(list, v)
			i++

			if i == max_return {
				return
			}
		}
	}

	return
}
func handler(w http.ResponseWriter, r *http.Request) {
	printParams(r)

	callback := r.FormValue("callback")
	wd := r.FormValue("wd")
	rs := r.FormValue("rs")
	if len(wd) < 2 {
		w.Write([]byte("[]"))
		return
	}

	if len(rs) < 2 {
		w.Write([]byte("[]"))
		return
	}

	haystack := haystacks[rs]
	if haystack == nil {
		w.Write([]byte("[]"))
		return
	}

	t := NewTick()
	list := search(haystack, []byte(wd))
	log.Printf("search %s in %d ms!\n", wd, t.Pick()/(1000*1000))

	ret := make([]interface{}, 0, len(list))
	for _, v := range list {
		ret = append(ret, v.Json())
	}

	b, _ := json.Marshal(ret)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if len(callback) == 0 {
		w.Write(b)
		return
	}
	w.Write([]byte(callback + "("))
	w.Write(b)
	w.Write([]byte(");"))
}

func Serv() {

	init_db()
	defer db.Close()

	loadHaystacks()

	http.HandleFunc("/res/", handler)
	http.ListenAndServe(":8888", nil)

}
