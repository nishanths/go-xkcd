package xkcd

import (
	"io"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	random = rand.New(
		&lockedRandSource{
			src: rand.NewSource(time.Now().UnixNano()),
		},
	)
}

// Client represents the HTTP client
// and any settings used to make requests
// to the xkcd API.
type Client struct {
	HTTPClient *http.Client
	Config
}

// NewClient returns a Client configured with sane default
// values.
func NewClient() *Client {
	return &Client{
		http.DefaultClient,
		Config{
			UseHTTPS: true,
		},
	}
}

func (c *Client) baseURL() string {
	protocol := "http://"

	if c.UseHTTPS {
		protocol = "https://"
	}

	return protocol + "xkcd.com"
}

// do performs a http request. If there is no error, the caller is responsible
// for closing the returned response body.
func (c *Client) do(req *http.Request) (io.ReadCloser, error) {
	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		res.Body.Close()
		return nil, newStatusError(res.StatusCode)
	}

	return res.Body, nil
}
