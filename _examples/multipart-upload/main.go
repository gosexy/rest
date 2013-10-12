// This is an example for the gosexy/rest package that issues a multipart POST
// request that contains an image to be uploaded and receives a JSON response
// that is directly converted to a map[string]interface{} value.

package main

import (
	"log"
	"menteslibres.net/gosexy/rest"
	"net/url"
	"os"
	"path"
)

func main() {

	var err error

	// You may want to see the full client debug.
	// rest.Debug = true

	filePath := "gopher.png"

	// First, check if the file to be uploaded is there.
	localFile, err := os.Open(filePath)

	if err != nil {
		log.Printf("Could not locate fle to upload %s.", err.Error())
	}

	// Destination variable.
	buf := map[string]interface{}{}

	// This service allows unauthenticated image uploads.
	requestURL := "http://cubeupload.com/upload_json.php"

	// These values must be passed in order for the file to be accepted.
	requestVariables := url.Values{
		"name":     {path.Base(filePath)},
		"userHash": {"false"},
		"userID":   {"false"},
	}

	// This is a file map that relates files with parameter names.
	fileMap := map[string][]rest.File{
		"fileinput[0]": []rest.File{
			rest.File{
				path.Base(localFile.Name()),
				localFile,
			},
		},
	}

	// The rest.NewMultipartBody creates a specially formatted body that mixes
	// parameters with encoded binary data.
	multipartBody, err := rest.NewMultipartBody(requestVariables, fileMap)

	if err != nil {
		log.Printf("Could not create multipart body %s.", err.Error())
	}

	// Let's pass buf's address ad first argument and issue a multipart POST request.
	err = rest.PostMultipart(&buf, requestURL, multipartBody)

	// Was there any error?
	if err == nil {

		// Printing response dump.
		log.Printf("Got response: buf = %v\n", buf)

		if name, ok := buf["file_name"].(string); ok {
			log.Printf("Your image was uploaded to http://cubeupload.com/im/%s", name)
		} else {
			log.Printf("Could not upload your image!")
		}

	} else {
		// Yes, we had an error.
		log.Printf(err.Error())
	}

}
