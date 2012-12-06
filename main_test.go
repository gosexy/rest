package rest

import (
	"encoding/json"
	//"fmt"
	"github.com/gosexy/sugar"
	"github.com/gosexy/to"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func init() {
	// Creating a new test server.
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			response := sugar.Tuple{
				"method": r.Method,
				"proto":  r.Proto,
				"host":   r.Host,
				"header": r.Header,
				"url":    r.URL.String(),
			}
			if r.Body != nil {
				response["body"], _ = ioutil.ReadAll(r.Body)
			}
			data, err := json.Marshal(response)
			if err == nil {
				w.Write(data)
			}
		},
	)
	go http.ListenAndServe("127.0.0.1:62621", nil)
}

func TestEnableDebug(t *testing.T) {
	Debug = true
}

func TestImplicitGet(t *testing.T) {
	client := New("http://www.golang.org")
	_, err := client.Text()
	if err != nil {
		t.Errorf("Failed test: %s\n", err.Error())
	}
}

func TestExplicitGet(t *testing.T) {

	var err error

	client := New("http://maps.googleapis.com/maps/api/")

	_, err = client.To("/geocode/").Do()

	if err != nil {
		t.Errorf("Failed test: %s\n", err.Error())
	}

	if client.Response.StatusCode != 404 {
		t.Errorf("Expecting 404, got %d\n", client.Response.StatusCode)
	}

	_, err = client.To("/geocode/json").Get(url.Values{
		"address": {"1600 Amphitheatre Parkway, Mountain View, CA"},
		"sensor":  {"true"},
	}).Text()

	if client.Response.StatusCode != 200 {
		t.Errorf("Expecting 200, got %d\n", client.Response.StatusCode)
	}

	_, err = client.To("/geocode/json").Post(url.Values{
		"address": {"1600 Amphitheatre Parkway, Mountain View, CA"},
		"sensor":  {"true"},
	}).Text()

	if client.Response.StatusCode != 200 {
		t.Errorf("Expecting 200, got %d\n", client.Response.StatusCode)
	}

	var data sugar.Tuple

	data, err = client.To("/geocode/json").Get(url.Values{
		"address": {"1600 Amphitheatre Parkway, Mountain View, CA"},
		"sensor":  {"true"},
	}).Json()

	if client.Response.StatusCode != 200 {
		t.Errorf("Expecting 200, got %d\n", client.Response.StatusCode)
	}

	if data.Get("status") != "OK" {
		t.Errorf("Failed test.")
	}

}

func TestRequestTypes(t *testing.T) {
	var data sugar.Tuple

	client := New("http://127.0.0.1:62621")

	client.To("/foo/bar").Do()

	if client.Response.StatusCode != 200 {
		t.Errorf("Expecting 200, got %d\n", client.Response.StatusCode)
	}

	data, _ = client.To("/foo/bar").Get(nil).Json()

	if client.Response.StatusCode != 200 {
		t.Errorf("Expecting 200, got %d\n", client.Response.StatusCode)
	}

	if data.Get("method") != "GET" {
		t.Errorf("Expecting GET, got %s\n", data.Get("method"))
	}

	data, _ = client.To("/foo/bar").Post(url.Values{"foo": {"bar"}}).Json()

	if data.Get("method") != "POST" {
		t.Errorf("Expecting POST, got %s\n", data.Get("method"))
	}

	data, _ = client.To("/foo/bar").Put(url.Values{"foo": {"bar"}}).Json()

	if data.Get("method") != "PUT" {
		t.Errorf("Expecting PUT, got %s\n", data.Get("method"))
	}

	data, _ = client.To("/foo/bar").Delete(url.Values{"foo": {"bar"}}).Json()

	if data.Get("method") != "DELETE" {
		t.Errorf("Expecting DELETE, got %s\n", data.Get("method"))
	}

	data, _ = client.To("/foo/bar").Head().Json()

	if client.Response.StatusCode != 200 {
		t.Errorf("Expecting 200, got %d\n", client.Response.StatusCode)
	}

	if data != nil {
		t.Errorf("Expecting HEAD, got %s\n", data.Get("method"))
	}

	data, _ = client.To("/foo/bar?a=b").Get(url.Values{"foo": {"bar"}}).Json()

	if client.Response.StatusCode != 200 {
		t.Errorf("Expecting 200, got %d\n", client.Response.StatusCode)
	}

	if data.Get("method") != "GET" {
		t.Errorf("Expecting GET, got %s\n", data.Get("method"))
	}

	if data.Get("url") != "/foo/bar?a=b&foo=bar" {
		t.Errorf("Expecting /foo/bar?a=b&foo=bar, got %s\n", data.Get("url"))
	}

}

func TestCustomHeader(t *testing.T) {
	var data sugar.Tuple

	client := New("http://127.0.0.1:62621")

	client.Header.Set("Foo", "Bar")
	client.Header.Set("User-Agent", "gosexy/rest")

	data, _ = client.To("/foo/bar?a=b").Get(url.Values{"foo": {"bar"}}).Json()

	if client.Response.StatusCode != 200 {
		t.Errorf("Expecting 200, got %d\n", client.Response.StatusCode)
	}

	if data.Get("method") != "GET" {
		t.Errorf("Expecting GET, got %s\n", data.Get("method"))
	}

	if data.Get("url") != "/foo/bar?a=b&foo=bar" {
		t.Errorf("Expecting /foo/bar?a=b&foo=bar, got %s\n", data.Get("url"))
	}

	if to.List(data.Get("header/Foo"))[0] != "Bar" {
		t.Errorf("Expecting Bar, got %s\n", to.List(data.Get("header/Foo"))[0])
	}

	if to.List(data.Get("header/User-Agent"))[0] != "gosexy/rest" {
		t.Errorf("Expecting gosexy/rest, got %s\n", to.List(data.Get("header/User-Agent"))[0])
	}

}
