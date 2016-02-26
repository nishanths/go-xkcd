package xkcd

import (
	"math"
	"testing"
)

func TestGetLatest(t *testing.T) {
	_, err := c.GetLatest()
	if err != nil {
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	comic, err := c.Get(614)

	if err != nil {
		t.Fail()
	}

	if comic.PublishDate.Month() != 7 {
		t.Fail()
	}
}

func TestGetNotFound(t *testing.T) {
	_, err := c.Get(math.MaxUint32)

	if err == nil {
		t.Fail()
	}

	xkcdErr, ok := err.(Error)

	// should be ok if the request succeeded
	// and we received a response
	if !ok {
		t.Fail()
	}

	if xkcdErr.StatusCode != 404 {
		t.Fail()
	}
}
