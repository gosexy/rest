/*
	A HTTP-REST client for Go that makes easy working with web services.
*/

package rest

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
)

const Version = "0.3"

/*
	Set true to log client requests and server responses to stdout.
*/
var Debug = false

var ioReadCloserType reflect.Type = reflect.TypeOf((*io.ReadCloser)(nil)).Elem()

type File struct {
	Name string
	io.Reader
}

type MultipartBody struct {
	contentType string
	buf         io.Reader
}

/*
	Client structure.
*/
type Client struct {
	Header http.Header
	Prefix string
}

/*
	Default client for requests.
*/
var DefaultClient = &Client{}

/*
	Creates a new client, all relative URLs this client receives will be prefixed
	by the given URL.
*/
func New(prefix string) (*Client, error) {
	var err error
	_, err = url.Parse(prefix)
	if err != nil {
		return nil, fmt.Errorf("Variable prefix must be a valid URL: %s", err.Error())
	}
	self := &Client{}
	self.Prefix = strings.TrimRight(prefix, "/") + "/"
	return self, nil
}

func (self *Client) newMultipartRequest(buf interface{}, method string, addr *url.URL, body *MultipartBody) error {
	var res *http.Response
	var req *http.Request

	var err error

	if body == nil {
		return fmt.Errorf("Could not create a multipart request without a body.")
	} else {
		req, err = http.NewRequest(
			method,
			addr.String(),
			body.buf,
		)
	}

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", body.contentType)

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

func (self *Client) newRequest(buf interface{}, method string, addr *url.URL, body *strings.Reader) error {
	var res *http.Response
	var req *http.Request

	var err error

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

/*
	Executes a PUT request and stores the response into the buf pointer.
*/
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

/*
	Executes a DELETE request and stores the response into the given buf pointer.
*/
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

/*
	Executes a multipart PUT request and stores the response into the given buf
	pointer.
*/
func (self *Client) PutMultipart(buf interface{}, uri string, data *MultipartBody) error {
	addr, err := url.Parse(self.Prefix + strings.TrimLeft(uri, "/"))

	if err != nil {
		return err
	}

	return self.newMultipartRequest(buf, "PUT", addr, data)
}

/*
	Executes a multipart POST request and stores the response into the given buf
	pointer.
*/
func (self *Client) PostMultipart(buf interface{}, uri string, data *MultipartBody) error {
	addr, err := url.Parse(self.Prefix + strings.TrimLeft(uri, "/"))

	if err != nil {
		return err
	}

	return self.newMultipartRequest(buf, "POST", addr, data)
}

/*
	Executes a POST request and stores the response into the given buf pointer.
*/
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

/*
	Executes a GET request and stores the response into the given buf pointer.
*/
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
	Creates a *MultipartBody based on the given params and map of files.
*/
func (self *Client) CreateMultipartBody(params url.Values, filemap map[string][]File) (*MultipartBody, error) {

	buf := bytes.NewBuffer(nil)

	body := multipart.NewWriter(buf)

	if filemap != nil {
		for key, files := range filemap {

			for _, file := range files {

				writer, err := body.CreateFormFile(key, path.Base(file.Name))

				if err != nil {
					return nil, err
				}

				_, err = io.Copy(writer, file.Reader)

				if err != nil {
					return nil, err
				}
			}
		}
	}

	if params != nil {
		for key, _ := range params {
			for _, value := range params[key] {
				body.WriteField(key, value)
			}
		}
	}

	body.Close()

	return &MultipartBody{body.FormDataContentType(), buf}, nil
}

/*
	Returns the body of the request as a io.ReadCloser
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

		log.Printf("\n")

	}

	return res, err
}

func Get(dest interface{}, uri string, data url.Values) error {
	return DefaultClient.Get(dest, uri, data)
}

func Post(dest interface{}, uri string, data url.Values) error {
	return DefaultClient.Post(dest, uri, data)
}

func Put(dest interface{}, uri string, data url.Values) error {
	return DefaultClient.Put(dest, uri, data)
}

func Delete(dest interface{}, uri string, data url.Values) error {
	return DefaultClient.Delete(dest, uri, data)
}

func PostMultipart(dest interface{}, uri string, data *MultipartBody) error {
	return DefaultClient.PostMultipart(dest, uri, data)
}

func PutMultipart(dest interface{}, uri string, data *MultipartBody) error {
	return DefaultClient.PutMultipart(dest, uri, data)
}
