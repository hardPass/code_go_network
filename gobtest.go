package main

import (
  "bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

var (
	gobfile = "d:\\test.gob"
)

func main() {
	//toFs()
	fromFs()
}

func encodeGob(obj interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func decodeGob(obj interface{}, n []byte) error {
	buf := bytes.NewBuffer(n)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(obj)

	return err
}

func toFs() {
	m := make(map[string]interface{})
	m["a"] = "aaa"
	m["b"] = "bbb"
	m["c"] = "ccc"

	n, _ := encodeGob(m)

	err := ioutil.WriteFile(gobfile, n, 0600)
	if err != nil {
		panic(err)
	}

}

func fromFs() {
	n, err := ioutil.ReadFile(gobfile)
	if err != nil {
		panic(err)
	}
	m := make(map[string]interface{})

	err = decodeGob(&m, n)
	if err != nil {
		panic(err)
	}

	fmt.Println(m)

}
