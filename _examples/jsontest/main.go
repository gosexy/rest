// This is an example for the gosexy/rest package that issues a POST request
// to a JSON service.

package main

import (
	// Import the gosexy/rest package.
	"github.com/gosexy/rest"
	"log"
	"net/url"
)

func main() {

	// You may want to see the full client debug.
	// rest.Debug = true

	// Destination variable.
	buf := map[string]interface{}{}

	// This service returns a JSON string containing your IP, like:
	// {"ip": "173.194.64.141"}
	requestURL := "http://ip.jsontest.com/"

	// We don't need any GET vars.
	requestVariables := url.Values{}

	// Let's pass buf's address ad first argument.
	err := rest.Get(&buf, requestURL, requestVariables)

	// Was there any error?
	if err == nil {

		// Printing response dump.
		log.Printf("Got response: buf = %v\n", buf)

		// Expecting a map with a single "ip" key.
		if ip, ok := buf["ip"].(string); ok {
			// What is my IP?
			log.Printf("According to ip.jsontest.com, your IP is %s\n", ip)
		}

	} else {
		// Yes, we had an error.
		log.Printf("Request failed: %s", err.Error())
	}

}
