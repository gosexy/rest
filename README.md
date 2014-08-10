# menteslibres.net/gosexy/rest

The `menteslibres.net/gosexy/rest` is an utility package that helps creating
HTTP-API clients with Go.

## Getting the package

```shell
go get -u menteslibres.net/gosexy/rest
```

## Capabilities

* GET, POST, PUT, DELETE.
* Automatic data conversion.
* HTTP basic authentication.
* Multipart requests.
* Multipart File uploads.
* Raw requests.
* Cookie jar.

## Special clients

## Debugging

Add `REST_DEBUG=1` to your list of enviroment variables to see all the talk
between client and server.

```
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
