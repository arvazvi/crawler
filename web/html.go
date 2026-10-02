package web

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"strings"

	queue "github.com/arvazvi/crawler/Queue"
	set "github.com/arvazvi/crawler/Set"
	"github.com/arvazvi/crawler/db"
	"golang.org/x/net/html"
)

type FetchOutput struct {
	Body []byte
	Error error
}

func GetHRef(t html.Token) (string, bool) {
	const HREF_TAG = "href"

	var (
		href string
		ok bool
	)
 
	for _, a := range t.Attr {
		if a.Key == HREF_TAG {
			if len(a.Val) == 0 || !strings.HasPrefix(a.Val, "http") {
				return a.Val, false
			}
			return a.Val, true
		}
		ok = true
		href = a.Val
	}
	return href, ok
}

func FetchPage(url string) ([]byte, error) {
	res, err := http.Get(url)
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

func DoFetch(url string, c chan *FetchOutput) {
	body, err := FetchPage(url)

	c <- &FetchOutput{
		Body: body,
		Error: err,
	}
}

type WebPage struct {
	URL 	string `json:"url"`
	Content string `json:"content"`
	Title 	string `json:"title"`
}

type Parser struct {
	Queue 	*queue.Queue
	Set 	*set.HashSet
	DB		*db.MongoDB
}

func (p *Parser) ParseHTML(url string, content []byte) (*WebPage, error) {
	const (
		MAX_CALLS = 1000
		MAX_TOKENS = 500
	)
	var (
		count int
		contentLength int
		hasBody bool
	)
	page := &WebPage{URL: url}
	tokenizer := html.NewTokenizer(bytes.NewReader(content))

	for {
		if tokenizer.Next() == html.ErrorToken || count > MAX_TOKENS {
			if p.Set.Size() < MAX_CALLS {
				log.Printf("Inserting %s", url)
				err := p.DB.Insert(context.TODO(), page)
				if err != nil {
					log.Printf("Err=%s", err.Error())
					return nil, err
				}
			}
			return page, nil
		}
		t := tokenizer.Token()
		if t.Type == html.StartTagToken {
			if t.Data == "body" {
				hasBody = true
			}
			if t.Data == "javascript" || t.Data == "script" || t.Data == "style" {
				// Skip script and style tags
				tokenizer.Next()
				continue
			}
			if t.Data == "title" {
				tokenizer.Next()
				title := tokenizer.Token().Data
				page.Title = title
				// log.Printf("Count: %d | %s -> %s\n", p.Set.Size(), url, title)
			}
			if t.Data == "a" {
				href, ok := GetHRef(t)
				if !ok {
					continue
				}
				if p.Set.Contains(href) || !strings.HasPrefix(href, "http") {
					continue
				} else {
					// log.Printf("Enqueing %s", href)
					p.Queue.Enqueue(href)
				}
			}
		}
		if hasBody && t.Type == html.TextToken && contentLength < MAX_TOKENS {
			page.Content += strings.TrimSpace(t.Data)
			contentLength += len(t.Data)
		}
		count++
	}
}

func NewParser(q *queue.Queue, s *set.HashSet, db *db.MongoDB) *Parser {
	return &Parser{Queue: q, Set: s, DB: db}
}