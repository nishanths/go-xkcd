package xkcd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const BaseURL = "https://xkcd.com"

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
}

func NewClient() *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    BaseURL,
	}
}

type comicResponse struct {
	Alt        string `json:"alt"`
	Day        string `json:"day"`
	Img        string `json:"img"`
	Link       string `json:"link"`
	Month      string `json:"month"`
	News       string `json:"news"`
	Num        int    `json:"num"`
	SafeTitle  string `json:"safe_title"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	Year       string `json:"year"`
}

type Comic struct {
	Alt        string
	Day        int
	ImageURL   string
	URL        string
	Month      int
	News       string
	Number     int
	SafeTitle  string
	Title      string
	Transcript string
	Year       int
}

type StatusError struct {
	Code int
}

var _ error = StatusError{}

func (e StatusError) Error() string {
	return fmt.Sprintf("bad response status code: %d", e.Code)
}

func (c *Client) Get(ctx context.Context, number int) (Comic, error) {
	return c.do(ctx, fmt.Sprintf("/%d/info.0.json", number))
}

func (c *Client) Latest(ctx context.Context) (Comic, error) {
	return c.do(ctx, fmt.Sprintf("/info.0.json"))
}

func (c *Client) do(ctx context.Context, reqPath string) (Comic, error) {
	req, err := http.NewRequest("GET", c.BaseURL+reqPath, nil)
	if err != nil {
		return Comic{}, fmt.Errorf("failed to build request: %v", err)
	}
	req = req.WithContext(ctx)

	rsp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Comic{}, fmt.Errorf("failed to do request: %v", err)
	}
	if rsp.StatusCode != 200 {
		return Comic{}, StatusError{Code: rsp.StatusCode}
	}
	defer drainAndClose(rsp.Body)

	var cr comicResponse
	if err := json.NewDecoder(rsp.Body).Decode(&cr); err != nil {
		return Comic{}, fmt.Errorf("failed to json-unmarshal response: %v", err)
	}

	d, err := strconv.Atoi(cr.Day)
	if err != nil {
		return Comic{}, fmt.Errorf("failed to parse day: %v", err)
	}
	m, err := strconv.Atoi(cr.Month)
	if err != nil {
		return Comic{}, fmt.Errorf("failed to parse month: %v", err)
	}
	y, err := strconv.Atoi(cr.Year)
	if err != nil {
		return Comic{}, fmt.Errorf("failed to parse year: %v", err)
	}

	return Comic{
		Alt:        cr.Alt,
		Day:        d,
		ImageURL:   cr.Img,
		URL:        cr.Link,
		Month:      m,
		News:       cr.News,
		Number:     cr.Num,
		SafeTitle:  cr.SafeTitle,
		Title:      cr.Title,
		Transcript: cr.Transcript,
		Year:       y,
	}, nil
}

func drainAndClose(r io.ReadCloser) {
	io.Copy(ioutil.Discard, r)
	r.Close()
}
