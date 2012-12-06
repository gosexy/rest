/*
	A HTTP-REST client for Go that makes easy working with
	web services.
*/

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

const Version = "0.1"

/*
	Set true to log client requests and server responses to stdout.
*/
var Debug = false

/*
	Client structure
*/
type Client struct {
	Request  *http.Request
	Response *http.Response
	Header   http.Header
	method   string
	prefix   string
	addr     *url.URL
	data     url.Values
	args     url.Values
	err      error
}

/*
	Creates a new client with a prefix URL.
*/
func New(prefix string) *Client {
	self := &Client{}
	self.Reset()
	self.SetPrefix(prefix)
	return self
}

/*
	Returns the latest error.
*/
func (self *Client) Error() error {
	return self.err
}

/*
	Sets the prefix for URLs.
*/
func (self *Client) SetPrefix(prefix string) {
	self.prefix = prefix
}

/*
	Resets the client for reusing it.
*/
func (self *Client) Reset() *Client {
	self.Request = nil
	self.Response = nil
	self.method = "GET"
	self.prefix = ""
	self.addr = nil
	self.data = url.Values{}
	self.args = url.Values{}
	self.err = nil
	self.Header = http.Header{
		"User-Agent": {fmt.Sprintf("gosexy/rest-%s", Version)},
	}
	return self
}

/*
	Sets the request URL. If prefix is not null, the given URL
	will be appended to the prefix.
*/
func (self *Client) To(addr string) *Client {
	var err error
	self.addr, err = url.Parse(strings.TrimRight(self.prefix, "/") + "/" + strings.TrimLeft(addr, "/"))
	if err != nil {
		self.err = err
	}
	self.args = nil
	return self
}

/*
	Prepares a HTTP HEAD request.
*/
func (self *Client) Head() *Client {
	self.method = "HEAD"
	return self
}

/*
	Prepares a HTTP PUT request.
*/
func (self *Client) Put(data url.Values) *Client {
	self.method = "PUT"
	self.data = data
	return self
}

/*
	Prepares a HTTP DELETE request.
*/
func (self *Client) Delete(data url.Values) *Client {
	self.method = "DELETE"
	self.data = data
	return self
}

/*
	Prepares a HTTP POST request.
*/
func (self *Client) Post(data url.Values) *Client {
	self.method = "POST"
	self.data = data
	return self
}

/*
	Prepares a HTTP GET request.
*/
func (self *Client) Get(data url.Values) *Client {
	self.args = data

	self.method = "GET"

	if self.addr != nil {
		if self.addr.RawQuery == "" {
			self.addr.RawQuery = self.args.Encode()
		} else {
			self.addr.RawQuery = self.addr.RawQuery + "&" + self.args.Encode()
		}
	} else {
		self.err = fmt.Errorf("No URL specified.")
	}

	return self
}

/*
	Executes a requests and returns a sugar.Tuple object created
	from the JSON response.
*/
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

/*
	Executes a requests and returns a text string of the body.
*/
func (self *Client) Text() (string, error) {
	bytes, err := self.Do()

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

/*
	Executes a request and returns the []byte array of the body.
*/
func (self *Client) Do() ([]byte, error) {
	var err error

	if self.addr == nil {
		self.To("")
	}

	switch self.method {
	case "POST", "PUT", "DELETE":
		self.Request, err = http.NewRequest(
			self.method,
			self.addr.String(),
			strings.NewReader(self.data.Encode()),
		)
		self.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	default:
		self.Request, err = http.NewRequest(
			self.method,
			self.addr.String(),
			nil,
		)
	}

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	for k, _ := range self.Header {
		self.Request.Header.Set(k, self.Header.Get(k))
	}

	self.Response, err = client.Do(self.Request)

	if err != nil {
		return nil, err
	}

	if self.Response.Header.Get("Content-Encoding") == "gzip" {
		self.Response.Body, err = gzip.NewReader(self.Response.Body)
		if err != nil {
			return nil, err
		}
	}

	defer self.Response.Body.Close()

	bytes, err := ioutil.ReadAll(self.Response.Body)

	if err != nil {
		return nil, err
	}

	if Debug == true {

		log.Printf("> %s %s", self.Request.Method, self.Request.Proto)
		for k, _ := range self.Request.Header {
			for kk, _ := range self.Request.Header[k] {
				log.Printf("> %s: %s", k, self.Request.Header[k][kk])
			}
		}

		log.Printf("< %s %s", self.Response.Proto, self.Response.Status)
		for k, _ := range self.Response.Header {
			for kk, _ := range self.Response.Header[k] {
				log.Printf("< %s: %s", k, self.Response.Header[k][kk])
			}
		}

		log.Printf("%s\n", string(bytes))
	}

	return bytes, nil
}
