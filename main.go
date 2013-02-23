/*
	A HTTP-REST client for Go that makes easy working with web services.
*/

package rest

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const Version = "0.3"

/*
	Set true to log client requests and server responses to stdout.
*/
var Debug = false

var ioReadCloserType reflect.Type = reflect.TypeOf((*io.ReadCloser)(nil)).Elem()

/*
	Client structure
*/
type Client struct {
	Header http.Header
	Prefix string
}

/*
	Creates a new client with a prefix URL.
*/
func New(prefix string) *Client {
	self := &Client{}
	self.Prefix = strings.TrimRight(prefix, "/") + "/"
	return self
}

func (self *Client) newRequest(buf interface{}, method string, addr *url.URL, body *strings.Reader) error {
	var res *http.Response
	var req *http.Request

	var err error

	fmt.Printf("READER: %v\n", body)
	fmt.Printf("URL: %v\n", addr.String())

	if body == nil {
		req, err = http.NewRequest(
			method,
			addr.String(),
			nil,
		)
	} else {
		req, err = http.NewRequest(
			method,
			addr.String(),
			body,
		)
	}

	switch method {
	case "POST", "PUT":
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}

	if err != nil {
		return err
	}

	res, err = self.Do(req)

	if err != nil {
		return err
	}

	err = self.handleResponse(buf, res)

	if err != nil {
		return err
	}

	return nil
}

func (self *Client) Put(buf interface{}, path string, data url.Values) error {
	var body *strings.Reader = nil

	addr, err := url.Parse(self.Prefix + strings.TrimLeft(path, "/"))

	if err != nil {
		return err
	}

	if data != nil {
		body = strings.NewReader(data.Encode())
	}

	return self.newRequest(buf, "PUT", addr, body)
}

func (self *Client) Delete(buf interface{}, path string, data url.Values) error {
	var body *strings.Reader = nil

	addr, err := url.Parse(self.Prefix + strings.TrimLeft(path, "/"))

	if err != nil {
		return err
	}

	if data != nil {
		body = strings.NewReader(data.Encode())
	}

	return self.newRequest(buf, "DELETE", addr, body)
}

func (self *Client) Post(buf interface{}, path string, data url.Values) error {
	var body *strings.Reader = nil

	addr, err := url.Parse(self.Prefix + strings.TrimLeft(path, "/"))

	if err != nil {
		return err
	}

	if data != nil {
		body = strings.NewReader(data.Encode())
	}

	return self.newRequest(buf, "POST", addr, body)
}

func (self *Client) Get(buf interface{}, path string, data url.Values) error {
	addr, err := url.Parse(self.Prefix + strings.TrimLeft(path, "/"))

	if err != nil {
		return err
	}

	if data != nil {
		if addr.RawQuery == "" {
			addr.RawQuery = data.Encode()
		} else {
			addr.RawQuery = addr.RawQuery + "&" + data.Encode()
		}
	}

	return self.newRequest(buf, "GET", addr, nil)

}

/*
func (self *Client) CreateMultipart(params url.Values, files map[string][]io.ReadCloser) {

	buf := bytes.NewBuffer(nil)
	body := multipart.NewWriter(buf)

	for key, file := range files {

		writer, err := body.CreateFormFile("media[]", path.Base(file))

		if err != nil {
			return nil, err
		}

		reader, err := os.Open(file)

		if err != nil {
			return nil, err
		}

		io.Copy(writer, reader)

		reader.Close()
	}

	params = merge(url.Values{"status": {status}}, params)

	//fullURI := Prefix + strings.Trim(endpoint, "/") + ".json"

	//self.client.SignParam(self.auth, "POST", fullURI, params)

	for k, _ := range params {
		for _, value := range params[k] {
			body.WriteField(k, value)
		}
	}

	body.Close()

	//fmt.Printf("%v\n", buf)

	req := &multipartBody{body.FormDataContentType(), buf}
}
*/

func (self *Client) Body(res *http.Response) (io.ReadCloser, error) {
	var body io.ReadCloser
	var err error

	if res.Header.Get("Content-Encoding") == "gzip" {
		body, err = gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
	} else {
		body = res.Body
	}

	return body, nil
}

func fromBytes(dst reflect.Value, buf []byte) error {
	var err error

	switch dst.Kind() {
	case reflect.String:
		// string
		dst.Set(reflect.ValueOf(string(buf)))
		return nil
	case reflect.Slice:
		switch dst.Type().Elem().Kind() {
		// []byte
		case reflect.Uint8:
			dst.Set(reflect.ValueOf(buf))
			return nil
		// []interface{}
		case reflect.Interface:
			t := []interface{}{}
			err = json.Unmarshal(buf, &t)

			if err == nil {
				dst.Set(reflect.ValueOf(t))
				return nil
			}
		}
	case reflect.Map:
		switch dst.Type().Elem().Kind() {
		case reflect.Interface:
			// map[string] interface{}
			m := map[string]interface{}{}

			err = json.Unmarshal(buf, &m)

			if err == nil {
				dst.Set(reflect.ValueOf(m))
				return nil
			}
		}
	}

	if err != nil {
		return err
	}

	return fmt.Errorf("Could not convert response (%s) to %s.", reflect.TypeOf(buf), dst.Type())
}

func (self *Client) handleResponse(dst interface{}, res *http.Response) error {
	body, err := self.Body(res)

	if err != nil {
		return err
	}

	if dst == nil {
		return nil
	}
	rv := reflect.ValueOf(dst)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("Destination is not a pointer.")
	}

	switch rv.Elem().Type() {
	case ioReadCloserType:
		rv.Elem().Set(reflect.ValueOf(body))
	default:
		buf, err := ioutil.ReadAll(body)

		if err != nil {
			return err
		}

		err = fromBytes(rv.Elem(), buf)

		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Client) Do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}

	// Copying headers
	for k, _ := range self.Header {
		req.Header.Set(k, self.Header.Get(k))
	}

	res, err := client.Do(req)

	if Debug == true {

		log.Printf("Fetching %v\n", req.URL.String())

		log.Printf("> %s %s", req.Method, req.Proto)
		for k, _ := range req.Header {
			for kk, _ := range req.Header[k] {
				log.Printf("> %s: %s", k, req.Header[k][kk])
			}
		}

		log.Printf("< %s %s", res.Proto, res.Status)
		for k, _ := range res.Header {
			for kk, _ := range res.Header[k] {
				log.Printf("< %s: %s", k, res.Header[k][kk])
			}
		}

	}

	return res, err
}
