package crawler

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"webcrawler/internal"
)

type Crawler struct {
	domain      string
	graph       SiteGraph
	maxDepth    int
	logger      *log.Logger
	workerCount int
}

func NewCrawler(domain string, count, depth int, loglevel string) Crawler {
	return Crawler{
		domain:      domain,
		graph:       NewSiteGraph(),
		maxDepth:    depth,
		workerCount: count,
		logger:      internal.NewLogger(loglevel),
	}
}

func (c *Crawler) Crawl() error {
	var wg sync.WaitGroup
	getWork, sendWork, resultCh := c.setupUrlPublisher()
	c.setupWorkerPool(&wg, sendWork, getWork, resultCh)
	wg.Wait()
	c.logger.Infof("Processed %d pages", len(c.graph.Graph))
	return nil
}

type urlInfo struct {
	url   string
	depth int
}

type workResult struct {
	url string
	err error
}

func (c *Crawler) setupUrlPublisher() (<-chan urlInfo, chan<- urlInfo, chan<- workResult) {
	sendWork := make(chan urlInfo)
	getWork := make(chan urlInfo)
	workProcessedCh := make(chan workResult)

	go func() {
		pendingWork := 1
		go func() {
			sendWork <- urlInfo{
				url:   c.domain,
				depth: 0,
			}
		}()
		for pendingWork > 0 {
			select {
			case result := <-workProcessedCh:
				if result.err != nil {
					c.logger.Warnf("Error processing url: %s. Error: %s", result.url, result.err)
				}
				pendingWork--
			case newWork := <-getWork:
				pendingWork++
				go func() {
					sendWork <- urlInfo{
						url:   c.domain + newWork.url,
						depth: newWork.depth,
					}
				}()
			}
		}
		close(sendWork)
	}()
	return sendWork, getWork, workProcessedCh
}

func (c *Crawler) setupWorkerPool(wwg *sync.WaitGroup, sendWork chan<- urlInfo, getWork <-chan urlInfo, resultCh chan<- workResult) {
	wwg.Add(c.workerCount)
	spawnWorker := func(id int) {
		defer wwg.Done()
		for work := range getWork {
			if visited := c.graph.VisitedPage(work.url); visited || work.depth >= c.maxDepth {
				resultCh <- workResult{}
				continue
			}
			c.logger.Debugf("Worker #%d. Processing url: %s \n", id+1, work.url)
			resp, err := http.Get(work.url)
			if err != nil {
				err = errors.Wrap(err, "http.Get")
				resultCh <- workResult{work.url, err}
				continue
			}

			title, links, err := ParseLinks(resp.Body, c.domain)
			resp.Body.Close()
			if err != nil {
				err = errors.Wrap(err, "ParseLinks")
				resultCh <- workResult{work.url, err}
				continue
			}

			c.graph.AddPage(work.url, Page{
				Title: title,
				Links: links,
			})

			for _, link := range links {
				sendWork <- urlInfo{
					url:   link,
					depth: work.depth + 1,
				}
			}
			resultCh <- workResult{work.url, nil}
		}
	}
	for i := 0; i < c.workerCount; i++ {
		go spawnWorker(i)
	}
}

func (c *Crawler) GetRelationships() map[string]Page {
	return c.graph.Graph
}
