package crawler

import (
	"bytes"
	"context"
	"log"
	"strings"

	"github.com/arvazvi/crawler/db"
	"github.com/arvazvi/crawler/queue"
	"github.com/arvazvi/crawler/set"
	"github.com/arvazvi/crawler/web"
	"golang.org/x/net/html"
)

type Fetch struct {
	Body []byte
	Error error
}

type Crawler struct {
	Queue 	*queue.Queue
	Crawled *set.HashSet
	DB		db.Database
	Channel chan *Fetch
}

type WebPage struct {
	URL 	string `json:"url"`
	Content string `json:"content"`
	Title 	string `json:"title"`
}

func NewCrawler(q *queue.Queue, s *set.HashSet, d db.Database) *Crawler {
	return &Crawler{Queue: q, Crawled: s, DB: d, Channel: make(chan *Fetch)}
}

func (c *Crawler) DoFetch(url string) {
	body, err := web.FetchPage(url)

	c.Channel <- &Fetch{
		Body: body,
		Error: err,
	}
}

func (c *Crawler) ParseHTML(ctx context.Context, url string, content []byte) {
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
			if c.Crawled.Size() < MAX_CALLS {
				err := c.DB.Insert(ctx, page)
				if err != nil {
					log.Printf("Err=%s", err.Error())
					return
				}
			}
			return
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
				page.Title = tokenizer.Token().Data
			}
			if t.Data == "a" {
				href, ok := web.GetHRef(t)
				if !ok {
					continue
				}
				if ok && c.Crawled.Contains(href) {
					// Already crawled
					continue
				} else {
					c.Queue.Enqueue(href)
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