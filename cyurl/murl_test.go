package cyurl

import (
	"log"
	"testing"
)

var data = [][]string{
	{"AA BB+", "AA%20BB%2B", "AA%20BB+"},
	{"AA&BB", "AA%26BB", "AA&BB"},
}

func TestQueryEscape(t *testing.T) {
	for _, d := range data {
		q := QueryEscape(d[0])
		if q != d[1] {
			log.Printf("%s != %s", q, d[1])
			t.Fail()
		}

		p := PathEscape(d[0])
		if p != d[2] {
			log.Printf("%s != %s", p, d[2])
			t.Fail()
		}

		qq, err := QueryUnescape(q)
		if err != nil {
			log.Printf("%s", err)
			t.Fail()
		}
		if qq != d[0] {
			log.Printf("%s != %s", qq, d[0])
			t.Fail()
		}

		pp, err := PathUnescape(p)
		if err != nil {
			log.Printf("%s", err)
			t.Fail()
		}
		if pp != d[0] {
			log.Printf("%s != %s", pp, d[0])
			t.Fail()
		}
	}
}
