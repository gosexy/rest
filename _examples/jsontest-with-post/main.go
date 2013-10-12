// This is an example for the gosexy/rest package that issues a POST request
// and receives a JSON string that is directly converted to a
// map[string]interface{} variable.

package main

import (
	"log"
	"menteslibres.net/gosexy/rest"
	"net/url"
)

func main() {

	// You may want to see the full client debug.
	// rest.Debug = true

	// Destination variable.
	buf := map[string]interface{}{}

	// This service returns a JSON string like:
	// {
	//	"md5": "fa4c6baa0812e5b5c80ed8885e55a8a6",
	//	"original": "example_text"
	// }
	requestURL := "http://md5.jsontest.com/"

	// We just need the "text" variable
	requestVariables := url.Values{
		"text": {"example_text"},
	}

	// Let's pass buf's address ad first argument and issue a POST request.
	err := rest.Post(&buf, requestURL, requestVariables)

	// Was there any error?
	if err == nil {

		// Printing response dump.
		log.Printf("Got response: buf = %v\n", buf)

		// Expecting a map with a single "md5" key.

		if hash, ok := buf["md5"].(string); ok {
			// What is my IP?
			log.Printf("According to md5.jsontest.com, the MD5 hash of %s is %s\n", requestVariables.Get("text"), hash)
		}

	} else {
		// Yes, we had an error.
		log.Printf(err.Error())
	}

}
