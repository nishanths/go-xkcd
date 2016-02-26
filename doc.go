/*
package xkcd provides methods to make requests to the xkcd.com API.

Example:

    client := xkcd.NewClient()
    comic, err := client.GetLatest()

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v", comic)


The Random() and RandomInRange() functions are safe to call concurrently.
They also do not modify the global random number generator.

More details on the xkcd API can be found at https://xkcd.com/json.html.
*/
package xkcd
