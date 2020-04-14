# go-xkcd

HTTP Client for the xkcd API.

[![wercker status](https://app.wercker.com/status/6c1de0bfd64a428d6ece5a2337268160/s "wercker status")](https://app.wercker.com/project/bykey/6c1de0bfd64a428d6ece5a2337268160) [![Coverage Status](https://coveralls.io/repos/github/nishanths/go-xkcd/badge.svg?branch=master)](https://coveralls.io/github/nishanths/go-xkcd?branch=master)
[![GoDoc](https://godoc.org/github.com/nishanths/go-xkcd?status.svg)](https://godoc.org/github.com/nishanths/go-xkcd)

[<img alt="https://xkcd.com/1481/" title="https://xkcd.com/1481/" src="http://imgs.xkcd.com/comics/api.png" width="250">](https://xkcd.com/1481/)

Details on the xkcd API can be found [here](https://xkcd.com/json.html).

## Import path

Import the package as:

```
github.com/nishanths/go-xkcd/v2
```

and refer to it as `xkcd`.

## Example

The following program prints details about [xkcd.com/599](http://xkcd.com/599):

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/nishanths/go-xkcd/v2"
)

func main() {
    client := xkcd.NewClient()
    
    comic, err := client.Get(context.Background(), 599)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s: %s\n", comic.Title, comic.ImageURL) // Apocalypse: http://imgs.xkcd.com/comics/apocalypse.png
}
```

## Test

To run tests:

```
$ go test -race
```

## Godoc

https://godoc.org/github.com/nishanths/go-xkcd

## License

The [MIT License](http://nishanths.mit-license.org). Copyright © [Nishanth Shanmugham](https://github.com/nishanths).
