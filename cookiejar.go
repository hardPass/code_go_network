package main

import (
	"code.google.com/p/mahonia"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

var (
	user_agent = "Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US) AppleWebKit/525.13 (KHTML, like Gecko) Chrome/0.2.149.29 Safari/525.13"
)

type Jar struct {
	sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.Lock()
	jar.cookies[u.Host] = cookies
	jar.Unlock()
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

func verifycode(client *http.Client) (string, error) {
	addr := "http://www.xxoo.com/codetest.asp"
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", user_agent)
	req.Header.Set("Referer", "http://www.xxoo.com/")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	cookies := client.Jar.Cookies(u)
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "verifycode" {
			return cookies[i].Value, nil
		}
	}

	return "", nil
}

func index(client *http.Client) {
	addr := "http://www.xxoo.com"
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", user_agent)
	req.Header.Set("Referer", "http://www.xxoo.com/")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()

}

func main() {
	proxy, err := url.Parse("http://127.0.0.1:8888") // for fiddler
	if err != nil {
		return
	}

	client := &http.Client{
		&http.Transport{
			Proxy: http.ProxyURL(proxy),
			// NO timeout
			// NO deadline
			// just a shit of test
		},
		nil,
		NewJar(),
	}

	// index(client)
	// client = &http.Client{nil, nil, NewJar()}
	code, err := verifycode(client)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("verify code : %s \n", code)

	addr := "http://www.xxoo.com/view11.asp?cllx=02&car_number=%B4%A8A597Q3&car_cj=443364&vcode=" + code + "&Submit=%B2%E9%D1%AF"
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", user_agent)
	req.Header.Set("Referer", "http://www.xxoo.com/")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	decoder := mahonia.NewDecoder("GBK")
	fmt.Printf("%s", decoder.ConvertString(string(b)))
}
