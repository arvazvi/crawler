package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/arvazvi/crawler/ai"
	"github.com/arvazvi/crawler/config"
	"github.com/arvazvi/crawler/crawler"
	"github.com/arvazvi/crawler/db"
	"github.com/arvazvi/crawler/queue"
	"github.com/arvazvi/crawler/set"
)

func main_() {
    const SEED = "https://www.cc.gatech.edu/"
    conf := config.Resolve()

    ctx := context.Background()
    db, err := db.NewConnection(ctx, conf)
    if err != nil {
        panic(err.Error())
    }

    ticker := time.NewTicker(10 * time.Second)
	done := make(chan bool)

    crawled := set.NewSet()
    q := queue.NewQueue()
    craw := crawler.NewCrawler(q, crawled, db)

    // Tick every minute
	go func() {
		for {
            select {
            case <-done:
                return
            }
        }
	}()

    craw.Queue.Enqueue(SEED)
	url, _ := q.Deque()
	craw.Crawled.Insert(url)

    go craw.DoFetch(url)

    result := <-craw.Channel
    if result.Error != nil {
        log.Printf("Error while processing url=%s", url)
    }
    craw.ParseHTML(ctx, url, result.Body)
    if err != nil {
        panic(err)
    }

    for q.Size() > 0 && crawled.Size() < 5000 {
		url, _ := q.Deque()

		crawled.Insert(url)
		go craw.DoFetch(url)

		result := <-craw.Channel
        if result.Error != nil {
            log.Printf("Error while processing url=%s", url)
            continue
        }
		craw.ParseHTML(ctx, url, result.Body)
        if err != nil {
            panic(err)
        }
	}

	ticker.Stop()
    done <- true
}
