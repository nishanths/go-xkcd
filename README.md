# [go-xkcd](https://github.com/nishanths/go-xkcd)

xkcd API client for Golang.

[![wercker status](https://app.wercker.com/status/6c1de0bfd64a428d6ece5a2337268160/s "wercker status")](https://app.wercker.com/project/bykey/6c1de0bfd64a428d6ece5a2337268160) [![Coverage Status](https://coveralls.io/repos/github/nishanths/go-xkcd/badge.svg?branch=master)](https://coveralls.io/github/nishanths/go-xkcd?branch=master)
[![GoDoc](https://godoc.org/github.com/nishanths/go-xkcd?status.svg)](https://godoc.org/github.com/nishanths/go-xkcd)

[<img alt="https://xkcd.com/1481/" title="https://xkcd.com/1481/" src="http://imgs.xkcd.com/comics/api.png" width="250">](https://xkcd.com/1481/)

Details on the xkcd API can be found [here](https://xkcd.com/json.html).

## Install

```
$ go get github.com/nishanths/go-xkcd
```

Import the package using `github.com/nishanths/go-xkcd` and refer to it as `xkcd`. 

Each major version has a separate branch. If you need a specific version, please clone the branch instead.

## Example

The following program prints details about [xkcd.com/599](http://xkcd.com/599):

```go
package main

import (
    "fmt"
    "log"

    "github.com/nishanths/go-xkcd"
)

func main() {
    client := xkcd.NewClient()
    comic, err := client.Get(599)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%s: %s", comic.Title, comic.ImageURL) // Apocalypse: http://imgs.xkcd.com/comics/apocalypse.png
}

```

## Test

To run tests:

```
$ go test
```

## Documentation

View Go Doc [online](https://godoc.org/github.com/nishanths/go-xkcd).

To view go docs locally, after installing the package run:

```
$ godoc -http=:6060
```

Then visit [`http://localhost:6060/pkg/github.com/nishanths/go-xkcd/`](http://localhost:6060/pkg/github.com/nishanths/go-xkcd/) in your browser.

#### Methods

The following methods are available on the client. All the methods return `(Comic, error)`.

* `Latest()`
* `Get(number int)`
* `Random()`
* `RandomInRange(begin, end, latest int)`

The fields available on `Comic` are:

```
type Comic struct {
	Alt         string
	PublishDate time.Time
	ImageURL    string
	URL         string
	News        string
	Number      int
	SafeTitle   string
	Title       string
	Transcript  string
}
```

## Contributing

Pull requests and issues are welcome!

To create a new pull request:

1. Fork the repository
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request on GitHub

## License

The [MIT License](http://nishanths.mit-license.org). Copyright © [Nishanth Shanmugham](https://github.com/nishanths).
