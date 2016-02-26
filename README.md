# go-xkcd

xkcd API client for Golang.

<img alt="https://xkcd.com/1481/" src="http://imgs.xkcd.com/comics/api.png" width="250">

Details on the xkcd API can be found [here](https://xkcd.com/json.html).

## Install

```
$ go get github.com/nishanths/go-xkcd
```

Import the package using `github.com/nishanths/go-xkcd` and refer to it as `xkcd`. 

Each major version has a separate branch. If you need a specific version, please clone the branch instead.

## Example

The following program prints details about the latest comic:

```go
package main

import (
    "fmt"
    "log"

    "github.com/nishanths/go-xkcd"
)

func main() {
    client := xkcd.NewClient()
    comic, err := client.GetLatest()

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%s: %s", comic.Title, comic.ImageURL)
}

```

## Test

To run tests:

```
$ go test
```

## Documentation

The following methods are available on the client. All the methods return `(Comic, error)`.

* `GetLatest()`
* `Get(number int)`
* `GetRandom(options ...int)` 

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