/*
Package xkcd provides a HTTP client for the xkcd.com JSON API.

Example:

    c := xkcd.NewClient()

    comic, err := c.Latest(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%+v\n", comic)


All methods on Client are safe to use concurrently.

More details on the xkcd API can be found at https://xkcd.com/json.html.
*/
package xkcd
