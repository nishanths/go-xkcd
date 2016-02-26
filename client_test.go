package xkcd

import (
	"net/http"
	"testing"
)

var c = &Client{
	http.DefaultClient,
	Config{
		UseHTTPS: true,
	},
}

func TestDo(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest("GET", "https://xkcd.com/info.0.json", nil)

	if err != nil {
		t.Fail()
	}

	_, err = c.do(req)

	if err != nil {
		t.Fail()
	}
}

func TestDoError(t *testing.T) {
	t.Parallel()
	req, _ := http.NewRequest("GET", "https://xkcd.com/F@Jfsf.json", nil)
	_, err := c.do(req)

	if err == nil {
		t.Fail()
	}
}
