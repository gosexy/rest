// This is an example for the gosexy/rest package.

package main

import (
	"github.com/gosexy/rest"
	"log"
	"net/url"
)

func main() {

	var err error

	// You may want to see the full client debug.
	// rest.Debug = true

	// Destination variable.
	buf := rest.Response{}

	// A nice gopher image.
	requestURL := "https://api.twitter.com/v1/foo.json"

	// We don't need any GET vars.
	requestVariables := url.Values{}

	// Let's pass buf's address ad first argument.
	err = rest.Get(&buf, requestURL, requestVariables)

	// Was there any error?
	if err == nil {

		// Printing response dump.
		log.Printf("Got response!")
		log.Printf("Response code: %d", buf.StatusCode)
		log.Printf("Response protocol version: %s", buf.Proto)
		log.Printf("Response length: %d", buf.ContentLength)
		log.Printf("Response header: %v", buf.Header)
		log.Printf("Response body: %s", string(buf.Body))
	} else {
		// Yes, we had an error.
		log.Printf("Request failed: %s", err.Error())
	}

}
