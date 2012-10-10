package rest

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gosexy/sugar"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var Debug = false

type Client struct {
	Request *http.Request
	method  string
	prefix  string
	addr    *url.URL
	data    *url.Values
	args    *url.Values
}

func New(prefix string) *Client {
	self := &Client{}
	self.Reset()
	self.prefix = prefix
	return self
}

func (self *Client) Reset() *Client {
	self.method = "GET"
	self.addr = &url.URL{}
	self.data = &url.Values{}
	self.args = &url.Values{}
	return self
}

func (self *Client) Post(data map[string]interface{}) *Client {
	self.method = "POST"
	self.data = &url.Values{}

	for key, _ := range data {
		self.data.Add(key, fmt.Sprintf("%v", data[key]))
	}

	return self
}

func (self *Client) Get(addr string, data map[string]interface{}) *Client {
	var err error

	if self.prefix != "" {
		self.addr, err = url.Parse(self.prefix + "/" + strings.Trim(addr, "/"))
	} else {
		self.addr, err = url.Parse(addr)
	}

	if err != nil {
		panic(err.Error())
	}

	self.args = &url.Values{}

	if data != nil {
		for key, _ := range data {
			self.args.Add(key, fmt.Sprintf("%v", data[key]))
		}
	}

	return self
}

func (self *Client) Json() (sugar.Tuple, error) {
	result := &sugar.Tuple{}

	bytes, err := self.Do()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, result)

	if err != nil {
		return nil, err
	}

	return *result, nil
}

func (self *Client) Text() (string, error) {
	bytes, err := self.Do()

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (self *Client) Do() ([]byte, error) {
	var req *http.Request
	var err error

	client := &http.Client{}

	if self.addr.RawQuery == "" {
		self.addr.RawQuery = self.args.Encode()
	} else {
		self.addr.RawQuery = self.addr.RawQuery + "&" + self.args.Encode()
	}

	switch self.method {
	case "POST", "PUT", "DELETE":
		req, err = http.NewRequest(self.method, self.addr.String(), strings.NewReader(self.data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	default:
		req, err = http.NewRequest(self.method, self.addr.String(), nil)
	}

	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.Header.Get("Content-Encoding") == "gzip" {
		res.Body, err = gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
	}

	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if Debug == true {
		log.Printf("%s %s", self.method, self.addr.String())
		log.Printf("DATA %v\n", self.data.Encode())
		log.Printf("RESPONSE %v\n", string(bytes))
		log.Printf("--\n")
	}

	self.Reset()

	return bytes, nil
}
