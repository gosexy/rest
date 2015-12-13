# menteslibres.net/gosexy/rest

The `rest` package helps you creating clients for HTTP APIs with Go.

## Getting the package

Get the `rest` package from `menteslibres.net/gosexy/rest` using `go get`:

```shell
go get -u menteslibres.net/gosexy/rest
```

## Usage and features

### Common HTTP verbs

The `rest` package comes with handy functions that are equivalent to HTTP
methods or *verbs*: `rest.Get()`, `rest.Post()`, `rest.Put()` and
`rest.Delete()`.

Let's take a look at the declaration of the `rest.Get()` function:

```go
func Get(dest interface{}, uri string, data url.Values) error {
  ...
}
```

as you can see `Get()` expects three variables: destination, url address and
parameters. The other *verb* functions expect the same variables.

The destination parameter must be either nil or a pointer to a variable. If you
provide a pointer, `rest` will try to do its best to convert the HTTP request's
body into the given type.

The second argument must be a fully qualified URL
(`http://www.example.com/foo/.../bar`) and the third argument must be either an
`url.Values{}` variable or nil if you don't need any parameter to be passed to
the URL.

In the following code example, the same request gets converted into different
types of destination variables:

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

If the server replies with a `Content-Type: application/json` header, it's
assumed that the response is a JSON formatted object, you can unmarshal JSON
objects if you provide a pointer to map or struct as destination.

This example expects JSON and unmarshals it to a map:

```go
buf := map[string]interface{}{}

// This service returns a JSON string containing your IP, like:
// {"ip": "173.194.64.141"}
rest.Get(&buf, "http://ip.jsontest.com", nil)

fmt.Printf("Got IP: %s", buf["ip"].(string))
```

This example expects JSON and unmarshals it into a struct:

```go
And also structs:

type ip_t struct {
  IP string `json:"ip"`
}

var buf ip_t

// This service returns a JSON string containing your IP, like:
// {"ip": "173.194.64.141"}
rest.Get(&buf, "http://ip.jsontest.com", nil)

fmt.Printf("Got IP: %s", buf.IP)
```

The `rest.Get()`, `rest.Post()`, `rest.Put()` and `rest.Delete()` functions,
use a default HTTP client, this is much like the `net/http`'s `DefaultClient`:

```go
var DefaultClient = new(Client)
```

### Custom clients

The `rest.Client` struct, allows you to create custom clients that use prefixes
that are to be automatically put at the beginning of the `url` argument in
`Get()`, `Post()`, `Put()` and `Delete()` methods.

```go
type Client struct {
  // These headers will be added in every request.
  Header http.Header
  // String to be added at the begining of every URL in Get(), Post(), Put()
  // and Delete() methods.
  Prefix string
  // Jar to store cookies.
  CookieJar *cookiejar.Jar
}
```

Using prefixes is useful to avoid repeating the first part of the URL if you're
using endpoints with similar names, for example:

```go
var customClient *rest.Client
var err error

if customClient, err = rest.New(`https://api.example.com/v1/`); err != nil {
  return err
}

// This call will prepend "https://api.example.com/v1/" to the given URL,
// so that the whole URL would be https://api.example.com/v1/users/add.
customClient.Get(&response, `/users/add`, url.Values{...})
```

You can also use custom clients to provide specific headers, you can set or get
those headers using the `Header` property of `rest.Client`.

```go
customClient.Header.Set(`X-Custom-Header`, `foo-api-version/v0.1`)
```

There is also a `CookieJar` property of the `rest.Client` type, this cookie jar
is created automatically and it stores the cookies that are received from the
site, if any.

### Basic authentication.

The `SetBasicAuth()` method of `rest.Client`, could be used to set required
information for basic authentication.

### Raw requests

The `PostRaw()` method of `rest.Client` allows you to post raw bytes to a given
URL.

```go
func (self *Client) PostRaw(dst interface{}, path string, body []byte) error {
  ...
}
```

This could be useful for some APIs that require you to post JSON-formatted
objects instead of plain old HTTP values.

### Multipart messages and file uploads

If you'd like to post a [multipart message][2], you can use the
`NewMultipartMessage()` function.

```go
message, err := NewMultipartMessage(url.Values{...}, nil)

if err != nil {
  ...
}

customClient.PostMultipart(&dst, "/api/multipart", message)
```

You can upload files using the multipart message encoding, as you can see in
the following example:

```go
fileToUpload, err := os.Open("my-avatar.png")

if err != nil {
  ...
}

defer fileToUpload.Close()

files := rest.FileMap{
  "variable_name": []rest.File{
    {
      Name: path.Base(fileToUpload.Name()),
      Reader: fileToUpload,
    },
  },
}

message, err := NewMultipartMessage(nil, files)

if err != nil {
  ...
}

customClient.PostMultipart(&dst, "/api/profile_photo/upload", message)
```

### Using detailed responses

`rest` provides an special type `rest.Response` that you can use when you need
to get detailed data from the response.

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
```

```go
buf := new(rest.Response)
rest.Get(&buf, "https://api.twitter.com/v1/foo.json", nil)
```

### Debugging

Add `REST_DEBUG=1` to your list of enviroment variables to see all the talk
between client and server.

```sh
REST_DEBUG=1 ./go-program
```

You can also use `rest.Debug()` to programmatically set the desired debug
level.

## Reference

See the [online docs][1] for `menteslibres.net/gosexy/rest` at [godoc.org][1].

## License

> Copyright (c) 2013-2014 José Carlos Nieto, https://menteslibres.net/xiam
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
[2]: http://www.w3.org/Protocols/rfc1341/7_2_Multipart.html
