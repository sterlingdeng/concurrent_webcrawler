package crawler

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync"
	"webcrawler/internal"
)

type Crawler struct {
	domain      string
	graph       SiteGraph
	logger      *log.Logger
	workerCount int
}

func NewCrawler(domain string, count int, loglevel string) Crawler {
	return Crawler{
		domain:      domain,
		graph:       NewSiteGraph(),
		workerCount: count,
		logger:      internal.NewLogger(loglevel),
	}
}

func (c *Crawler) Crawl() error {
	var wg sync.WaitGroup
	getWork, sendWork, resultCh := c.setupUrlPublisher(c.domain)
	c.setupWorkerPool(&wg, sendWork, getWork, resultCh)
	wg.Wait()
	c.logger.Infof("Processed %d pages", len(c.graph.Graph))
	return nil
}

type workResult struct {
	url string
	err error
}

func (c *Crawler) setupUrlPublisher(entry string) (<-chan string, chan<- string, chan<- workResult) {
	sendWork := make(chan string)
	getWork := make(chan string)
	workProcessedCh := make(chan workResult)

	go func() {
		pendingWork := 1
		go func() {
			sendWork <- entry
		}()
		for pendingWork > 0 {
			select {
			case result := <-workProcessedCh:
				if result.err != nil {
					c.logger.Warnf("Error processing url: %s. Error: %s", result.url, result.err)
				}
				pendingWork--
			case url := <-getWork:
				pendingWork++
				go func() { sendWork <- c.domain + url }()
			}
		}
		close(sendWork)
	}()
	return sendWork, getWork, workProcessedCh
}

func (c *Crawler) setupWorkerPool(wwg *sync.WaitGroup, sendWork chan<- string, getWork <-chan string, resultCh chan<- workResult) {
	wwg.Add(c.workerCount)
	spawnWorker := func(id int) {
		defer wwg.Done()
		for url := range getWork {
			if visited := c.graph.VisitedPage(url); visited {
				resultCh <- workResult{}
				continue
			}
			c.logger.Debugf("Worker #%d. Processing url: %s \n", id+1, url)
			resp, err := http.Get(url)
			if err != nil {
				err = errors.Wrap(err, "http.Get")
				resultCh <- workResult{url, err}
				continue
			}
			body, ok := resp.Body.(io.Reader)
			if !ok {
				err = errors.New("resp.Body type assertion to io.Reader failed")
				resultCh <- workResult{url, err}
				continue
			}

			title, links, err := ParseLinks(body, c.domain)
			resp.Body.Close()
			if err != nil {
				err = errors.Wrap(err, "ParseLinks")
				resultCh <- workResult{url, err}
				continue
			}

			c.graph.AddPage(url, Page{
				Title: title,
				Links: links,
			})

			for _, link := range links {
				sendWork <- link
			}
			resultCh <- workResult{url, nil}
		}
	}
	for i := 0; i < c.workerCount; i++ {
		go spawnWorker(i)
	}
}

func (c *Crawler) GetRelationships() map[string]Page {
	return c.graph.Graph
}
