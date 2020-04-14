package xkcd

import (
	"context"
	"flag"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	responseLatest = []byte(`{"month":"4","num":2293,"link":"","year":"2020","news":"","safe_title":"RIP John Conway","transcript":"","alt":"1937-2020","img":"https://imgs.xkcd.com/comics/rip_john_conway.gif","title":"RIP John Conway","day":"13"}`)
	response144    = []byte(`{"month":"8","num":144,"link":"","year":"2006","news":"","safe_title":"Parody Week: A Softer World","transcript":"when we open the lab each morning, we tell the robot to kill\nit's our little joke\nbut secretly\nwe're just afraid\nto tell it to love","alt":"The robot is pregnant.  It isn't mine.","img":"https://imgs.xkcd.com/comics/a_softer_robot.jpg","title":"Parody Week: A Softer World","day":"17"}`)

	comicLatest = Comic{
		Month:     4,
		Number:    2293,
		Year:      2020,
		SafeTitle: "RIP John Conway",
		Alt:       "1937-2020",
		ImageURL:  "https://imgs.xkcd.com/comics/rip_john_conway.gif",
		Title:     "RIP John Conway",
		Day:       13,
	}
	comic144 = Comic{
		Month:      8,
		Number:     144,
		Year:       2006,
		SafeTitle:  "Parody Week: A Softer World",
		Transcript: "when we open the lab each morning, we tell the robot to kill\nit's our little joke\nbut secretly\nwe're just afraid\nto tell it to love",
		Alt:        "The robot is pregnant.  It isn't mine.",
		ImageURL:   "https://imgs.xkcd.com/comics/a_softer_robot.jpg",
		Title:      "Parody Week: A Softer World",
		Day:        17,
	}
)

const nonExistentComic = math.MaxUint32

type testClient struct {
	desc   string
	client *Client
}

var (
	localClient, liveClient *Client
	clients                 []testClient
)

func localServer() *httptest.Server {
	writeJson := func(b []byte) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write(b)
		})
	}

	mux := http.NewServeMux()
	mux.Handle("/info.0.json", writeJson(responseLatest))
	mux.Handle("/144/info.0.json", writeJson(response144))
	mux.HandleFunc(fmt.Sprintf("/%d/info.0.json", nonExistentComic), func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	return httptest.NewUnstartedServer(mux)
}

func TestMain(m *testing.M) {
	flag.Parse()

	svr := localServer()
	svr.Start()
	defer svr.Close()

	localClient = &Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    svr.URL,
	}
	liveClient = NewClient()

	clients = []testClient{
		{"local", localClient},
		{"live", liveClient},
	}

	os.Exit(m.Run())
}

func TestLatest(t *testing.T) {
	t.Parallel()

	comic, err := localClient.Latest(context.Background())
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if comicLatest != comic {
		t.Errorf("comics do not match:\nexpected=%+v\n     got=%+v", comicLatest, comic)
		return
	}

	_, err = liveClient.Latest(context.Background())
	if err != nil {
		t.Errorf("%v", err)
		return
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	for _, c := range clients {
		t.Run(c.desc, func(t *testing.T) {
			comic, err := c.client.Get(context.Background(), 144)
			if err != nil {
				t.Errorf("%v", err)
				return
			}
			if comic144 != comic {
				t.Errorf("comics do not match:\nexpected=%+v\n     got=%+v", comic144, comic)
				return
			}
		})
	}
}

func TestStatusError(t *testing.T) {
	t.Parallel()

	for _, c := range clients {
		t.Run(c.desc, func(t *testing.T) {
			_, err := c.client.Get(context.Background(), math.MaxUint32)
			if err == nil {
				t.Errorf("error unexpectedly nil")
				return
			}
			statusErr, ok := err.(StatusError)
			if !ok {
				t.Errorf("error unexpectedly nil")
				return
			}
			if statusErr.Code != 404 {
				t.Errorf("wrong status code:\nexpected=%d\n     got=%d", 404, statusErr.Code)
				return
			}
		})
	}
}
