// This is an example for the gosexy/rest package that downloads a JPEG image
// and opens it with the image/jpeg package.

package main

import (
	"bytes"
	"github.com/gosexy/rest"
	"image/jpeg"
	"log"
	"net/url"
)

func main() {

	var err error

	// You may want to see the full client debug.
	// rest.Debug = true

	// Destination variable (an empty bytes.Buffer).
	buf := bytes.NewBuffer(nil)

	// A nice gopher image.
	requestURL := "http://talks.golang.org/2012/splash/appenginegophercolor.jpg"

	// We don't need any GET vars.
	requestVariables := url.Values{}

	// Let's pass buf's address ad first argument.
	err = rest.Get(&buf, requestURL, requestVariables)

	// Was there any error?
	if err == nil {

		// Printing some response data.
		log.Printf("Got response with size %d\n", buf.Len())

		// Trying to open the buffer with image/jpeg.
		log.Printf("Trying to decode JPEG file.\n")

		img, err := jpeg.Decode(buf)

		if err == nil {
			log.Printf("JPEG decoded correctly!")
			log.Printf("-> bounds: %v", img.Bounds())
		} else {
			log.Printf("Error decoding PNG file: %s", err.Error())
		}

	} else {
		// Yes, we had an error.
		log.Printf("Request failed: %s", err.Error())
	}

}
