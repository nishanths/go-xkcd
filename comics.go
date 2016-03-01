package xkcd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type comicResponse struct {
	Alt        string `json:"alt"`
	Day        string `json:"day"`
	ImageURL   string `json:"img"`
	URL        string `json:"link"`
	Month      string `json:"month"`
	News       string `json:"news"`
	Number     int    `json:"num"`
	SafeTitle  string `json:"safe_title"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	Year       string `json:"year"`
}

// Comic represents information and metadata
// about a single xkcd comic.
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

// UnmarshalJSON unmarshals the response from the xkcd enpoint
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
	comic.ImageURL = aux.ImageURL
	comic.URL = aux.URL
	comic.Number = aux.Number
	comic.News = aux.News
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

// Latest returns the latest comic's information
func (c *Client) Latest() (Comic, error) {
	return c.doComicRequest("/info.0.json")
}

// Get returns the comic for the specified number.
// The number is the comic number in the xkcd.com url.
// For example: https://xkcd.com/193.
func (c *Client) Get(number int) (Comic, error) {
	numStr := strconv.Itoa(number)
	return c.doComicRequest("/" + numStr + "/info.0.json")
}

func (c *Client) getLatestComicNumber() (int, error) {
	comic, err := c.Latest()
	if err != nil {
		return 0, err
	}
	return comic.Number, nil
}

// Random returns a random comic.
// The underlying random number generator's behavior may not match
// the behavior of the Random button on xkcd.com.
// Random never performs a request for a non-existent comic number.
//
// Also see: RandomInRange()
func (c *Client) Random() (Comic, error) {
	return c.RandomInRange(-1, -1, -1)
}

// RandomInRange returns a random comic using the given options.
// The underlying random number generator's behavior may not match
// the behavior of the Random button on xkcd.com.
//
// [begin, end) specify the range that the randomly chosen comic number
// can be in. If  begin equals -1, begin defaults to
// the number of the first comic. Likewise, if end equals
// -1, end defaults to the number of the latest comic + 1.
//
// latest specfies the number of the latest xkcd comic. Specifying the
// number eliminates the overhead of performing an additional HTTP request
// to find this number. Pass in -1 to let RandomInRange() find the latest
// comic number by performing the additional request.
func (c *Client) RandomInRange(begin, end, latest int) (comic Comic, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	if begin == -1 {
		begin = 1 // 1 is the first xkcd comic
	}

	if latest == -1 {
		latest, err = c.getLatestComicNumber()
		if err != nil {
			return
		}
	}

	if end == -1 {
		end = latest + 1
	}

	number := randomInt(begin, end)
	comic, err = c.Get(number)
	return
}
