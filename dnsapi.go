package main

import (
  "bytes"
	"fmt"
	"io/ioutil"
	pkglog "log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	login_email    = "************"
	login_password = "******"
	format         = "json"
	domain_id      = "s"
	record_id      = "s"
	sub_domain     = "www"
	record_line    = "默认"
	current_ip     = ""
	log            = pkglog.New(os.Stdout, "", pkglog.Ldate|pkglog.Ltime)
)

const (
	timeout  = 10 * time.Second
	interval = 30 * time.Second
)

func get_public_ip() (string, error) {
	conn, err := net.DialTimeout("tcp", "ns1.dnspod.net:6666", timeout)
	defer func() {
		if x := recover(); x != nil {
			log.Println("Can't get public ip", x)
		}
		if conn != nil {
			conn.Close()
		}
	}()
	if err == nil {
		var bytes []byte
		deadline := time.Now().Add(timeout)
		err = conn.SetDeadline(deadline)
		if err != nil {
			return "", err
		}
		bytes, err = ioutil.ReadAll(conn)
		if err == nil {
			return string(bytes), nil
		}
	}
	return "", err
}

func timeoutDialler(timeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		c, err := net.DialTimeout(netw, addr, timeout)
		if err != nil {
			return nil, err
		}
		deadline := time.Now().Add(timeout)
		err = c.SetDeadline(deadline)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
}

func domainList() bool {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: timeoutDialler(timeout),
		},
	}

	body := url.Values{
		"login_email":    {login_email},
		"login_password": {login_password},
		"format":         {format},
		"type":           {"all"},
	}
	req, err := http.NewRequest("POST", "https://dnsapi.cn/Domain.List", strings.NewReader(body.Encode()))
	req.Header.Set("Accept", "text/json")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return false

	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(bytes))
	return resp.StatusCode == 200
}

func recordList() bool {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: timeoutDialler(timeout),
		},
	}

	body := url.Values{
		"login_email":    {login_email},
		"login_password": {login_password},
		"format":         {format},
		"domain_id":      {domain_id},
	}
	req, err := http.NewRequest("POST", "https://dnsapi.cn/Record.List", strings.NewReader(body.Encode()))
	req.Header.Set("Accept", "text/json")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return false

	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(bytes))
	return resp.StatusCode == 200
}

func update_dnspod(ip string) bool {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: timeoutDialler(timeout),
		},
	}

	body := url.Values{
		"login_email":    {login_email},
		"login_password": {login_password},
		"format":         {format},
		"domain_id":      {domain_id},
		"record_id":      {record_id},
		"sub_domain":     {sub_domain},
		"record_line":    {record_line},
		"value":          {ip},
	}
	req, err := http.NewRequest("POST", "https://dnsapi.cn/Record.Ddns", strings.NewReader(body.Encode()))
	req.Header.Set("Accept", "text/json")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return false

	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(bytes))
	return resp.StatusCode == 200
}

func main() {
	go func() {
		for {
			ip, err := get_public_ip()
			if ip != "" && err == nil {
				log.Println("got ip:" + ip + " dns_ip is:" + current_ip)
				if ip != current_ip {
					log.Println("update dnspod with new ip:" + ip)
					if update_dnspod(ip) {
						current_ip = ip
					}
				}
			} else {
				log.Println("error:", err)
			}
			time.Sleep(interval)
		}
	}()

	http.HandleFunc("/", handler)          // redirect all urls to the handler function
	err := http.ListenAndServe(":80", nil) // listen for connections at port 9999 on the local machine
	log.Println(err)
}

var (
	list = NewList(17)
)

type syncLimitedList struct {
	sync.Mutex
	counter int
	list    []string
}

func NewList(cap int) *syncLimitedList {
	return &syncLimitedList{counter: 0, list: make([]string, cap)}
}

func (l *syncLimitedList) put(s string) {
	l.Lock()
	defer l.Unlock()

	l.list[l.counter%len(l.list)] = s
	l.counter++
}

func (l *syncLimitedList) desc() string {
	l.Lock()
	defer l.Unlock()

	times := 0
	i := l.counter % len(l.list)
	buffer := bytes.NewBuffer(make([]byte, 0, 1000))
	for {
		i--
		if i < 0 {
			i = len(l.list) - 1
		}
		if l.list[i] == "" {
			break
		}

		times++
		if times >= len(l.list) {
			break
		}

		buffer.WriteString("<p>")
		buffer.WriteString(l.list[i])
		buffer.WriteString("</p>")
	}

	return buffer.String()
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		return
	}
	line := time.Now().Format("2006-01-02 15:04:05") + "\t" + r.RemoteAddr
	log.Println(r.RemoteAddr)
	list.put(line)
	fmt.Fprintf(w, list.desc())
}
