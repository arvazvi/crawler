package web

import (
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func GetHRef(t html.Token) (href string, ok bool) {
	for _, a := range t.Attr {
        if a.Key == "href" {
			if len(a.Val) == 0 || !strings.HasPrefix(a.Val, "http") {
				ok = false
	            href = a.Val
				return href, ok
			}
            href = a.Val
			ok = true
        }
    }
    return href, ok
}

func FetchPage(url string) ([]byte, error) {
	res, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
