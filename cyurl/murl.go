package cyurl

import (
	"net/url"
	"strings"
)

func QueryEscape(v string) string {
	return strings.ReplaceAll(url.QueryEscape(v), "+", "%20")
}

func QueryUnescape(v string) (string, error) {
	return url.QueryUnescape(v)
}

func PathEscape(v string) string {
	return url.PathEscape(v)
}

func PathUnescape(v string) (string, error) {
	return url.PathUnescape(v)
}
