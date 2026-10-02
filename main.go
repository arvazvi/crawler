package main

import (
	"context"
	"fmt"
	"log"
	"time"

	queue "github.com/arvazvi/crawler/Queue"
	set "github.com/arvazvi/crawler/Set"
	"github.com/arvazvi/crawler/config"
	"github.com/arvazvi/crawler/db"
	"github.com/arvazvi/crawler/web"
)

type CrawlerStats struct {
	pagesPerMinute string // 0 0 \n 1 100
	crawledRatioPerMinute string 
	startTime time.Time
}

func (c *CrawlerStats) update(crawled *set.HashSet, queue *queue.Queue, t time.Time) {
	c.pagesPerMinute += fmt.Sprintf("%f %d\n", t.Sub(c.startTime).Minutes(), crawled.Size())
	c.crawledRatioPerMinute += fmt.Sprintf("%f %f\n", t.Sub(c.startTime).Minutes(), float64(crawled.Size())/float64(queue.Size()))
}

func (c *CrawlerStats) print() {
	fmt.Println("Pages crawled per minute:")
	fmt.Println(c.pagesPerMinute)
	fmt.Println("Crawl to Queued Ratio per minute:")
	fmt.Println(c.crawledRatioPerMinute)
}

func main() {
    const SEED = "https://www.cc.gatech.edu/"
    conf := config.Resolve()

    ctx := context.Background()
    db, err := db.NewConnection(ctx, conf)
    if err != nil {
        panic(err.Error())
    }

    ticker := time.NewTicker(1 * time.Minute)
	done := make(chan bool)
	crawlerStats := CrawlerStats{pagesPerMinute: "0 0\n", crawledRatioPerMinute: "0 0\n", startTime: time.Now()}

    crawled := set.NewSet()
    q := queue.NewQueue()
    parser := web.NewParser(q, crawled, db)

    // Tick every minute
	go func() {
		for {
            select {
            case <-done:
                return
            case t := <-ticker.C:
                crawlerStats.update(crawled, q, t)
            }
        }
	}()

    q.Enqueue(SEED)
	url, _ := q.Deque()
	crawled.Insert(url)
	c := make(chan *web.FetchOutput)

    go web.DoFetch(url, c)

    result := <-c
    if result.Error != nil {
        log.Fatalf("Error while processing url=%s", url)
    }
    _, err = parser.ParseHTML(url, result.Body)
    if err != nil {
        panic(err)
    }

    for q.Size() > 0 && crawled.Size() < 5000 {
		url, _ := q.Deque()

		crawled.Insert(url)
		go web.DoFetch(url, c)

		result := <-c
        if result.Error != nil {
            fmt.Println(result.Error)
            log.Fatalf("Error while processing url=%s", url)
        }
		_, err := parser.ParseHTML(url, result.Body)
        if err != nil {
            panic(err)
        }
	}

	ticker.Stop()
    done <- true
	fmt.Println("\n------------------CRAWLER STATS------------------")
	fmt.Printf("Total queued: %d\n", q.Total())
	fmt.Printf("To be crawled (Queue) size: %d\n", q.Size())
	fmt.Printf("Crawled size: %d\n", crawled.Size())
	crawlerStats.print()
}