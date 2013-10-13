# gosexy/rest

The `gosexy/rest` package adds some functionality on top of `net/http` in
order to make working with web services even easier.

## Quick tour

Whenever possible, `gosexy/rest` tries to convert the HTTP response into the
datatype you need, for example:

```go
// Dumping request response as a byte array.
bytesBuf := []byte{}
rest.Get(&bytesBuf, "http://golang.org", nil)

// Dumping request response as a string.
stringBuf := ""
rest.Get(&stringBuf, "http://golang.org", nil)

// Dumping request response into a bytes.Buffer buffer.
buf := bytes.NewBuffer(nil)
rest.Get(&buf, "http://golang.org", nil)
```

It does also support JSON (for JSON formatted documents).

```go
buf := map[string]interface{}

// This service returns a JSON string containing your IP, like:
// {"ip": "173.194.64.141"}
rest.Get(&buf, "http://ip.jsontest.com", nil)

fmt.Printf("Got IP: %s", buf["ip"].(string))
```

And if you need the whole document with complete headers and response code, a
`rest.Response` type is also provided:

```go
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

...

buf := rest.Response{}
rest.Get(&buf, "https://api.twitter.com/v1/foo.json", nil)
```

## Getting gosexy/rest

You can install this package as usual:

```
go get -u menteslibres.net/gosexy/rest
```

## Examples

This is a full code example, you can see more examples in the
[examples](./_examples) directory.

```go
// This is an example for the gosexy/rest package that issues a POST request
// to a JSON service.

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
```

## Documentation

See the [online docs][1] for `gosexy/rest` at [godoc.org][1].

## License

> Copyright (c) 2013 José Carlos Nieto, https://menteslibres.net/xiam
>
> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:
>
> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
> LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
> OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
> WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

[1]: http://godoc.org/menteslibres.net/gosexy/rest
