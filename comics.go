package xkcd

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type comicResponse struct {
	Alt        string `json:"alt"`
	Day        string `json"day"`
	Image      string `json:"img"`
	Link       string `json:"link"`
	Month      string `json:"month"`
	Number     int    `json:"num"`
	SafeTitle  string `json:"safe_title"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	Year       string `json:"year"`
}

type Comic struct {
	Alt         string
	PublishDate time.Time
	Image       string
	Link        string
	Month       string
	Number      int
	SafeTitle   string
	Title       string
	Transcript  string
	Year        string
}

func (comic *Comic) UnmarshalJSON(data []byte) error {
	var aux comicResponse
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	date, err := time.Parse("2006/1/2", aux.Year+"/"+aux.Month+"/"+aux.Day)
	if err != nil {
		return err
	}

	comic.Alt = aux.Alt
	comic.PublishDate = date
	comic.Image = aux.Image
	comic.Link = aux.Link
	comic.Number = aux.Number
	comic.SafeTitle = aux.SafeTitle
	comic.Title = aux.Title
	comic.Transcript = aux.Transcript

	return nil
}

func (c *Client) doComicRequest(path string) (Comic, error) {
	var comic Comic

	req, err := http.NewRequest("GET", c.baseURL()+path, nil)
	if err != nil {
		return comic, err
	}

	body, err := c.do(req)
	if err != nil {
		return comic, err
	}
	defer body.Close()

	if err := json.NewDecoder(body).Decode(&comic); err != nil {
		return comic, err
	}
	return comic, nil
}

func (c *Client) GetLatest() (Comic, error) {
	return c.doComicRequest("/info.0.json")
}

func (c *Client) Get(number int) (Comic, error) {
	numStr := strconv.Itoa(number)
	return c.doComicRequest("/" + numStr + "/info.0.json")
}
