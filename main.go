/*
  Copyright (c) 2013 José Carlos Nieto, https://menteslibres.net/xiam

  Permission is hereby granted, free of charge, to any person obtaining
  a copy of this software and associated documentation files (the
  "Software"), to deal in the Software without restriction, including
  without limitation the rights to use, copy, modify, merge, publish,
  distribute, sublicense, and/or sell copies of the Software, and to
  permit persons to whom the Software is furnished to do so, subject to
  the following conditions:

  The above copyright notice and this permission notice shall be
  included in all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
  MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
  NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
  LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
  OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
  WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// Utility package for Go that makes working with web services even easier.

package rest

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"reflect"
	"strings"
)

const Version = "0.3"

// Set true to log client requests and server responses to stdout.
var Debug = false

var (
	ioReadCloserType reflect.Type = reflect.TypeOf((*io.ReadCloser)(nil)).Elem()
	bytesBufferType  reflect.Type = reflect.TypeOf((**bytes.Buffer)(nil)).Elem()
	restResponseType reflect.Type = reflect.TypeOf((*Response)(nil)).Elem()
)

var (
	ErrInvalidPrefix           = errors.New(`URL prefix is invalid: %s`)
	ErrCouldNotCreateMultipart = errors.New(`Couldn't create a multipart request without a body.`)
	ErrCouldNotConvert         = errors.New(`Could not convert response %s to %s.`)
	ErrDestinationNotAPointer  = errors.New(`Destination is not a pointer.`)
)

// Can be used as a response value, useful when you need to work with headers,
// status codes and such.
type Response struct {
	Status        string
	StatusCode    int
	Proto         string
	ProtoMajor    int
	ProtoMinor    int
	ContentLength int64
	http.Header
	Body []byte
}

// This type can be used to represent a file that you'll later upload within
// a multipart request.
type File struct {
	Name string
	io.Reader
}

// Multipart body for multipart requests, you can't generate a MultipartBody
// directly, use rest.NewMultipartBody() instead.
type MultipartBody struct {
	contentType string
	buf         io.Reader
}

// A client, useful in case you need to communicate with an API and you'd like
// to use the same prefix for all of your requests or in scenarios where it
// would be handy to keep a session cookie.
type Client struct {
	Header    http.Header
	Prefix    string
	CookieJar *cookiejar.Jar
}

// Default client
var DefaultClient = &Client{}

// Creates a new client, relative URLs will be prefixed with the given prefix
// value
func New(prefix string) (*Client, error) {
	var err error
	_, err = url.Parse(prefix)
	if err != nil {
		return nil, fmt.Errorf(ErrInvalidPrefix.Error(), err.Error())
	}
	self := &Client{}
	self.Prefix = strings.TrimRight(prefix, "/") + "/"
	self.Header = http.Header{}
	self.CookieJar, err = cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return self, nil
}

func (self *Client) newMultipartRequest(buf interface{}, method string, addr *url.URL, body *MultipartBody) error {
	var res *http.Response
	var req *http.Request

	var err error

	if body == nil {
		return ErrCouldNotCreateMultipart
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

	res, err = self.do(req)

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
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		}
	}

	if err != nil {
		return err
	}

	res, err = self.do(req)

	if err != nil {
		return err
	}

	err = self.handleResponse(buf, res)

	if err != nil {
		return err
	}

	return nil
}

// Performs a HTTP PUT request, stores response into the buf pointer.
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

// Performs a HTTP DELETE request, stores response into the buf pointer.
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

// Performs a multipart HTTP PUT request, stores response into the buf pointer.
func (self *Client) PutMultipart(buf interface{}, uri string, data *MultipartBody) error {
	addr, err := url.Parse(self.Prefix + strings.TrimLeft(uri, "/"))

	if err != nil {
		return err
	}

	return self.newMultipartRequest(buf, "PUT", addr, data)
}

// Performs a multipart HTTP POST request, stores response into the buf pointer.
func (self *Client) PostMultipart(buf interface{}, uri string, data *MultipartBody) error {
	addr, err := url.Parse(self.Prefix + strings.TrimLeft(uri, "/"))

	if err != nil {
		return err
	}

	return self.newMultipartRequest(buf, "POST", addr, data)
}

// Performs a HTTP POST request, stores response into the buf pointer.
func (self *Client) PostRaw(buf interface{}, path string, data []byte) error {
	var body *strings.Reader = nil

	addr, err := url.Parse(self.Prefix + strings.TrimLeft(path, "/"))

	if err != nil {
		return err
	}

	if data != nil {
		body = strings.NewReader(string(data))
	}

	return self.newRequest(buf, "POST", addr, body)
}

// Performs a HTTP POST request, stores response into the buf pointer.
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

// Performs a HTTP GET request, stores response into the buf pointer.
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

// Creates a *MultipartBody based on the given params and map of files.
func NewMultipartBody(params url.Values, filemap map[string][]File) (*MultipartBody, error) {

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

// Returns the body of the request as a io.ReadCloser
func (self *Client) body(res *http.Response) (io.ReadCloser, error) {
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

	return fmt.Errorf(ErrCouldNotConvert.Error(), reflect.TypeOf(buf), dst.Type())
}

func (self *Client) handleResponse(dst interface{}, res *http.Response) error {

	body, err := self.body(res)

	if err != nil {
		return err
	}

	if dst == nil {
		return nil
	}
	rv := reflect.ValueOf(dst)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrDestinationNotAPointer
	}

	t := res.Header.Get("Content-Type")

	switch rv.Elem().Type() {
	case restResponseType:
		var err error

		r := Response{}

		r.Body, err = ioutil.ReadAll(body)

		if err != nil {
			return err
		}

		r.Header = res.Header
		r.Status = res.Status
		r.StatusCode = res.StatusCode
		r.Proto = res.Proto
		r.ProtoMajor = res.ProtoMajor
		r.ProtoMinor = res.ProtoMinor
		r.ContentLength = res.ContentLength

		rv.Elem().Set(reflect.ValueOf(r))
	case ioReadCloserType:
		rv.Elem().Set(reflect.ValueOf(body))
	case bytesBufferType:
		buf, err := ioutil.ReadAll(body)

		if err != nil {
			return err
		}

		dst := bytes.NewBuffer(buf)

		rv.Elem().Set(reflect.ValueOf(dst))
	default:
		buf, err := ioutil.ReadAll(body)

		if err != nil {
			return err
		}

		if strings.HasPrefix(t, "application/json") == true {
			if rv.Elem().Kind() == reflect.Struct || rv.Elem().Kind() == reflect.Map {
				err = json.Unmarshal(buf, dst)
				return err
			}
		}

		err = fromBytes(rv.Elem(), buf)

		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Client) do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}

	// Adding cookie jar
	if self.CookieJar != nil {
		client.Jar = self.CookieJar
	}

	// Copying headers
	for k, _ := range self.Header {
		req.Header.Set(k, self.Header.Get(k))
	}

	if req.Body == nil {
		req.Header.Del("Content-Type")
		req.Header.Del("Content-Length")
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

// Performs a HTTP GET request using the default client. Stores response
// in dest pointer.
func Get(dest interface{}, uri string, data url.Values) error {
	return DefaultClient.Get(dest, uri, data)
}

// Performs a HTTP POST request using the default client. Stores response
// in dest pointer.
func Post(dest interface{}, uri string, data url.Values) error {
	return DefaultClient.Post(dest, uri, data)
}

// Performs a HTTP PUT request using the default client. Stores response
// in dest pointer.
func Put(dest interface{}, uri string, data url.Values) error {
	return DefaultClient.Put(dest, uri, data)
}

// Performs a HTTP DELETE request using the default client. Stores response
// in dest pointer.
func Delete(dest interface{}, uri string, data url.Values) error {
	return DefaultClient.Delete(dest, uri, data)
}

// Performs a multipart HTTP POST request using the default client. Stores
// response in dest pointer.
func PostMultipart(dest interface{}, uri string, data *MultipartBody) error {
	return DefaultClient.PostMultipart(dest, uri, data)
}

// Performs a multipart HTTP PUT request using the default client. Stores
// response in dest pointer.
func PutMultipart(dest interface{}, uri string, data *MultipartBody) error {
	return DefaultClient.PutMultipart(dest, uri, data)
}
