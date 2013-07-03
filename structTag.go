package main

import (
  "fmt"
	"reflect"
	"strings"
)

type Address struct {
	City, State string "nima"
	Name        string `xml:"per son " json:"sjson" `
}

func main() {

	output := Address{"cc", "xx", "nmmm"}
	v := reflect.ValueOf(output)
	tt := reflect.TypeOf(reflect.Indirect(v).Interface())
	// tt := reflect.TypeOf(reflect.Indirect(v))

	for i := 0; i < tt.NumField(); i++ {
		fmt.Println("-----------------------------------------")
		fmt.Println(tt.Field(i).Name)
		fmt.Println(tt.Field(i).Tag)
		fmt.Println(tt.Field(i).Tag.Get("xml"))
		fmt.Println(tt.Field(i).Tag.Get("json"))

		// if bb.Get("beedb") == "PK" || reflect.ValueOf(bb).String() == "PK" {
		// 	orm.PrimaryKey = tt.Field(i).Name
		// }
	}

}
