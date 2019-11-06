package crawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

type Crawler struct {
	domain string
	graph  SiteGraph

	workerCount int
}

func NewCrawler(domain string, count int) Crawler {
	return Crawler{
		domain:      domain,
		graph:       NewSiteGraph(),
		workerCount: count,
	}
}

func (c *Crawler) Crawl() {
	err := c.crawl()
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Crawler) crawl() error {
	var pending int32
	urlChan := make(chan string)

	done := c.createWorkers(urlChan, c.graph, &pending)
	atomic.AddInt32(&pending, 1)
	urlChan <- c.domain
	<-done
	return nil
}

func (c *Crawler) createWorkers(urlChan chan string, sg SiteGraph, counter *int32) <-chan struct{} {
	var wg sync.WaitGroup
	wg.Add(c.workerCount)

	counterTick := make(chan struct{})
	decrementCounter := func() {
		atomic.AddInt32(counter, -1)
		counterTick <- struct{}{}
	}

	for i := 0; i < c.workerCount; i++ {
		go func(id int) {
			defer wg.Done()
			for url := range urlChan {
				if visited := sg.VisitedPage(url); visited {
					go decrementCounter()
					continue
				}
				fmt.Printf("Worker #%d. Processing url: %s \n", id, url)
				resp, err := http.Get(url)
				if err != nil {
					go decrementCounter()
					continue
				}
				body, ok := resp.Body.(io.Reader)
				if !ok {
					go decrementCounter()
					continue
				}

				title, links, err := ParseLinks(body, c.domain)
				resp.Body.Close()
				if err != nil {
					go decrementCounter()
					continue
				}

				sg.AddPage(url, Page{
					Title: title,
					Links: links,
				})

				for _, link := range links {
					atomic.AddInt32(counter, 1)
					go func(l string) {
						urlChan <- c.domain + l
					}(link)
				}
				go decrementCounter()
			}
		}(i)
	}

	doneCh := make(chan struct{})
	go func() {
		for *counter > 0 {
			<-counterTick
		}
		close(urlChan)
		wg.Wait()
		doneCh <- struct{}{}
	}()
	return doneCh
}

func (c *Crawler) GetRelationships() map[string]Page {
	return c.graph.Graph
}
