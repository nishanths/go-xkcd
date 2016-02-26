/*
package xkcd provides methods to make requests to the xkcd.com API. Details on the xkcd API can be found at https://xkcd.com/json.html.

Example:

    client := xkcd.NewClient()
    comic, err := client.GetLatest()

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v", comic)
*/
package xkcd
