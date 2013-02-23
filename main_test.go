package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

const testServer = "127.0.0.1:62621"
const reqForm = 1024 * 1024 * 8

var client *Client

func init() {
	// Creating a new test server.
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseMultipartForm(reqForm)

			getValues, _ := url.ParseQuery(r.URL.RawQuery)

			response := map[string]interface{}{
				"method": r.Method,
				"proto":  r.Proto,
				"host":   r.Host,
				"header": r.Header,
				"url":    r.URL.String(),
				"get":    getValues,
				"post":   r.Form,
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
	go http.ListenAndServe(testServer, nil)

	time.Sleep(time.Second * 1)
}

func TestEnableDebug(t *testing.T) {
	Debug = true
}

func TestInit(t *testing.T) {
	client = New("http://" + testServer)
}

func TestGet(t *testing.T) {
	var buf map[string]interface{}
	var err error

	err = client.Get(&buf, "/search", url.Values{"term": {"some string"}})

	if err != nil {
		t.Errorf("Failed test: %s\n", err.Error())
	}

	if buf["method"].(string) != "GET" {
		t.Errorf("Test failed.")
	}

	if buf["url"].(string) != "/search?term=some+string" {
		t.Errorf("Test failed.")
	}

	if buf["get"].(map[string]interface{})["term"].([]interface{})[0].(string) != "some string" {
		t.Errorf("Test failed.")
	}

	err = client.Get(&buf, "/search", nil)

	if err != nil {
		t.Errorf("Failed test: %s\n", err.Error())
	}

	if buf["method"].(string) != "GET" {
		t.Errorf("Test failed.")
	}
}

func TestPost(t *testing.T) {
	var buf map[string]interface{}
	var err error

	err = client.Post(&buf, "/search?foo=the+quick", url.Values{"bar": {"brown fox"}})

	if err != nil {
		t.Errorf("Failed test: %s\n", err.Error())
	}

	fmt.Printf("%v\n", buf)

	if buf["method"].(string) != "POST" {
		t.Errorf("Test failed.")
	}

	if buf["post"].(map[string]interface{})["bar"].([]interface{})[0].(string) != "brown fox" {
		t.Errorf("Test failed.")
	}

	if buf["get"].(map[string]interface{})["foo"].([]interface{})[0].(string) != "the quick" {
		t.Errorf("Test failed.")
	}

	err = client.Post(&buf, "/search?foo=the+quick", nil)

	if err != nil {
		t.Errorf("Failed test: %s\n", err.Error())
	}

	if buf["method"].(string) != "POST" {
		t.Errorf("Test failed.")
	}

	if buf["get"].(map[string]interface{})["foo"].([]interface{})[0].(string) != "the quick" {
		t.Errorf("Test failed.")
	}
}

func TestPut(t *testing.T) {
	var buf map[string]interface{}
	var err error

	err = client.Put(&buf, "/search?foo=the+quick", url.Values{"bar": {"brown fox"}})

	if err != nil {
		t.Errorf("Failed test: %s\n", err.Error())
	}

	fmt.Printf("%v\n", buf)

	if buf["method"].(string) != "PUT" {
		t.Errorf("Test failed.")
	}

	if buf["post"].(map[string]interface{})["bar"].([]interface{})[0].(string) != "brown fox" {
		t.Errorf("Test failed.")
	}

	if buf["get"].(map[string]interface{})["foo"].([]interface{})[0].(string) != "the quick" {
		t.Errorf("Test failed.")
	}

	err = client.Put(&buf, "/search?foo=the+quick", nil)

	if err != nil {
		t.Errorf("Failed test: %s\n", err.Error())
	}

	if buf["method"].(string) != "PUT" {
		t.Errorf("Test failed.")
	}

	if buf["get"].(map[string]interface{})["foo"].([]interface{})[0].(string) != "the quick" {
		t.Errorf("Test failed.")
	}
}

func TestDelete(t *testing.T) {
	var buf map[string]interface{}
	var err error

	err = client.Delete(&buf, "/search?foo=the+quick", url.Values{"bar": {"brown fox"}})

	if err != nil {
		t.Errorf("Failed test: %s\n", err.Error())
	}

	fmt.Printf("%v\n", buf)

	if buf["method"].(string) != "DELETE" {
		t.Errorf("Test failed.")
	}

	if buf["get"].(map[string]interface{})["foo"].([]interface{})[0].(string) != "the quick" {
		t.Errorf("Test failed.")
	}

	err = client.Delete(&buf, "/search?foo=the+quick", nil)

	if err != nil {
		t.Errorf("Failed test: %s\n", err.Error())
	}

	if buf["method"].(string) != "DELETE" {
		t.Errorf("Test failed.")
	}

	if buf["get"].(map[string]interface{})["foo"].([]interface{})[0].(string) != "the quick" {
		t.Errorf("Test failed.")
	}
}
