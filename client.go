package xkcd

import (
	"io"
	"net/http"
)

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

// do performs a http request. If the request was successful the
// response body and a nil error are returned. If the request failed or
// there was an error in the response, a nil error is returned.
// The returned response body has to be closed by the caller.
func (c *Client) do(req *http.Request) (io.ReadCloser, error) {
	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		return res.Body, nil
	}

	return nil, NewError(res.StatusCode)
}
