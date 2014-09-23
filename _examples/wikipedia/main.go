// This is an example for the gosexy/rest package that reads a Wikipedia page
// and counts the ocurrences of the word "Go".

package main

import (
	"github.com/gosexy/rest"
	"log"
	"net/url"
	"strings"
)

func main() {

	// You may want to see the full client debug.
	// rest.Debug = true

	// Destination variable.
	buf := ""

	// The Wikipedia article on Go.
	requestURL := "http://en.wikipedia.org/wiki/Golang"

	// We don't need any GET vars.
	requestVariables := url.Values{}

	// Let's pass buf's address ad first argument.
	err := rest.Get(&buf, requestURL, requestVariables)

	// Was there any error?
	if err == nil {

		// Printing response.
		log.Printf("Got response with length %d\n", len(buf))

		// How many times does the word "Go" appear in this page?
		log.Printf("The word \"Go\" appears %d times within the document.\n", strings.Count(buf, "Go"))

	} else {
		// Yes, we had an error.
		log.Printf("Request failed: %s", err.Error())
	}

}
