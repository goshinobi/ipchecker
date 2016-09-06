package ipchecker

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	buffer []func(*http.Client) string
	p      int
)

func init() {
	rand.Seed(time.Now().UnixNano())
	RegistChecker(inetIP)
	RegistChecker(httpbin)
	RegistChecker(cman)
	p = rand.Int() % len(buffer)
}

func RegistChecker(f func(*http.Client) string) {
	buffer = append(buffer, f)
}

func Check(c *http.Client) string {
	result := buffer[p](c)
	p++
	if p >= len(buffer) {
		p = 0
	}
	return result
}

var httpbin = func(c *http.Client) string {
	url := "https://httpbin.org/ip"
	resp, err := c.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	var ip = map[string]string{}
	if err = json.Unmarshal(byteArray, &ip); err != nil {
		return ""
	}
	return ip["origin"]
}

var inetIP = func(c *http.Client) string {
	url := "http://inet-ip.info/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "curl")

	resp, err := c.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	result := string(byteArray)
	result = strings.Replace(result, "\n", "", -1)
	result = strings.Replace(result, " ", "", -1)
	return result
}

var cman = func(c *http.Client) string {
	url := "https://www.cman.jp/network/support/go_access.cgi"

	resp, err := c.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return ""
	}

	return doc.Find(".outIp").Text()
}
