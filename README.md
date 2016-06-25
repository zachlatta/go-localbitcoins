# go-localbitcoins

go-localbitcoins is a Go wrapper around the [LocalBitcoins API](https://localbitcoins.com/api-docs).

**Documentation**: http://godoc.org/github.com/zachlatta/go-localbitcoins/localbitcoins

**Build Status**: [![Build Status](https://travis-ci.org/zachlatta/go-localbitcoins.png?branch=master)](https://travis-ci.org/zachlatta/go-localbitcoins)

**Test Coverage**: [![Coverage Status](https://coveralls.io/repos/zachlatta/go-localbitcoins/badge.png?branch=master)](https://coveralls.io/r/zachlatta/go-localbitcoins?branch=master)

go-localbitcoins requires Go version 1.1 or greater.

## Usage

```go
import "github.com/zachlatta/go-localbitcoins/localbitcoins"
```

Construct a new LocalBitcoins client, then use services to access different parts of the API. For example, to get the user info for the account "zrl":

```go
client := localbitcoins.NewClient(nil)
zrl, _, err := client.Accounts.Get("zrl")
```

For complete usage, see the full [package docs](http://godoc.org/github.com/zachlatta/go-localbitcoins/localbitcoins).

### Authentication

go-localbitcoins does not directly handle authentication. Instead, when creating a new client, pass in an `http.Client` that can handle authentication for you. The easiest way to do this is by using the [goauth2-localbitcoins](https://github.com/zachlatta/goauth2-localbitcoins) fork of [goauth2](https://code.google.com/p/goauth2/) modified to play nice with LocalBitcoins. Further details regarding authentication on LocalBitcoins are available at https://localbitcoins.com/api-docs/#toc1.

A complete example with authentication is available at https://github.com/zachlatta/go-localbitcoins/blob/master/examples/example.go

## Acknowledgements

go-localbitcoins is heavily inspired by the wonderful [go-github](https://github.com/google/go-github) library.

## License

The MIT License (MIT)

Copyright (c) 2014 Zach Latta

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
