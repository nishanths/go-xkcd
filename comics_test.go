package xkcd

import (
	"fmt"
	"math"
	"sync"
	"testing"
)

type Option struct {
	begin, end, latest int
}

func TestLatest(t *testing.T) {
	t.Parallel()

	_, err := c.Latest()
	if err != nil {
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	comic, err := c.Get(614)

	if err != nil {
		t.Fail()
	}

	if comic.PublishDate.Month() != 7 {
		t.Fail()
	}
}

func TestGetNotFound(t *testing.T) {
	t.Parallel()

	_, err := c.Get(math.MaxUint32)

	if err == nil {
		t.Fail()
	}

	statusErr, ok := err.(StatusError)

	// should be ok if the request succeeded
	// and we received a response
	if !ok {
		t.Fail()
	}

	if statusErr.StatusCode != 404 {
		t.Fail()
	}
}

func TestRandomInRange(t *testing.T) {
	t.Parallel()
	var wg sync.WaitGroup

	options := []Option{
		Option{100, 101, 9999}, // end is the real limiting factor
		Option{-1, 2, 9999},
		Option{10, 1000, 9999},
		Option{1, -1, 1}, // limited by the one comic that has been published
	}
	wg.Add(len(options))

	for _, r := range options {
		go func(r Option) {
			defer wg.Done()
			comic, err := c.RandomInRange(r.begin, r.end, r.latest)

			if err != nil {
				t.Fail()
			}

			upperLimExcl := r.end
			if r.end == -1 {
				upperLimExcl = r.latest + 1
			}

			if comic.Number < r.begin || comic.Number >= upperLimExcl {
				fmt.Println(r)
				t.Fail()
			}
		}(r)
	}

	wg.Wait()
}

func TestRandomInRangeError(t *testing.T) {
	t.Parallel()
	var wg sync.WaitGroup

	options := []Option{
		Option{0, 0, -1},
		Option{100, 5, 9999},
	}
	wg.Add(len(options))

	for _, r := range options {
		go func(r Option) {
			defer wg.Done()
			_, err := c.RandomInRange(r.begin, r.end, r.latest)

			if err == nil {
				t.Fail()
			}
		}(r)
	}

	wg.Wait()
}

func TestRandom(t *testing.T) {
	t.Parallel()

	_, err := c.Random()
	if err != nil {
		t.Fail()
	}
}
