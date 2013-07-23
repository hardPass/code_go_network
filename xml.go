package main

import (
  "bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	confg_xml = flag.String("c", "config.xml", "set config xml")
)

type Receiver struct {
	PathIn   string `xml:"Path>In"`
	PathTmp  string `xml:"Path>Tmp"`
	Interval int64  `xml:"Interval"`
	Max      int64  `xml:"Max"`
}

type Service struct {
	Name       string    `xml:"name,attr"`
	Class      string    `xml:"class,attr"`
	Concurrent int64     `xml:"Concurrent"`
	Receiver   *Receiver `xml:"Receiver"`
}

func (svc *Service) String() string {
	buffer := bytes.NewBuffer(make([]byte, 0, 100))
	buffer.WriteString("Name:")
	buffer.WriteString(svc.Name)
	buffer.WriteString("\n")

	buffer.WriteString("Class:")
	buffer.WriteString(svc.Class)
	buffer.WriteString("\n")

	buffer.WriteString("Concurrent:")
	buffer.WriteString(strconv.Itoa(int(svc.Concurrent)))
	buffer.WriteString("\n")

	buffer.WriteString("\tPathIn:")
	buffer.WriteString(svc.Receiver.PathIn)
	buffer.WriteString("\n")

	buffer.WriteString("\tPathTmp:")
	buffer.WriteString(svc.Receiver.PathTmp)
	buffer.WriteString("\n")

	buffer.WriteString("\tInterval:")
	buffer.WriteString(strconv.Itoa(int(svc.Receiver.Interval)))
	buffer.WriteString("\n")

	buffer.WriteString("\tMax:")
	buffer.WriteString(strconv.Itoa(int(svc.Receiver.Max)))
	buffer.WriteString("\n")

	return buffer.String()
}

func main() {
	flag.Parse()

	xmlFile, err := os.Open(*confg_xml)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	total := 0
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		// Inspect the type of the token just read.
		switch tk := t.(type) {
		case xml.StartElement:
			if tk.Name.Local == "Service" {
				total++
				var svc *Service
				for _, attr := range tk.Attr {
					fmt.Printf("%s : %s \n", attr.Name, attr.Value)
				}

				decoder.DecodeElement(&svc, &tk)

				fmt.Println(svc.String())

			}
		default:
		}

	}
	fmt.Printf("Total services: %d \n", total)

}
