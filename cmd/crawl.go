package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
	"time"
	"webcrawler/internal/crawler"
)

const (
	DefaultDomain = "http://localhost:8080"
	DefaultWorkerCount = 10
)

var (
	DomainUrl string
	WorkerCount int
)

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawls a given website and returns the sitemap",
	Long:  "Crawls a given website, if none is provided, it will crawl monzo.com",
	Run: func(cmd *cobra.Command, args []string) {
		Crawl()
	},
}

func init() {
	crawlCmd.Flags().StringVarP(&DomainUrl, "domain", "d", DefaultDomain, "Domain for webcrawler to crawl")
	crawlCmd.Flags().IntVarP(&WorkerCount, "workers", "w", DefaultWorkerCount, "Sets the number of concurrent workers")
	rootCmd.AddCommand(crawlCmd)
}

func Crawl() {

	go func() {
		tickCh := time.Tick(time.Second)
		for {
			<-tickCh
			fmt.Println(runtime.NumGoroutine())
		}
	}()

	t1 := time.Now()
	c := crawler.NewCrawler(DomainUrl, WorkerCount)
	c.Crawl()
	graph := c.GetRelationships()
	for url, page := range graph {
		fmt.Printf("URL: %s Page Data: %+v \n", url, page)
	}
	t2 := time.Now()
	fmt.Println("time:", t2.Sub(t1))
	for true {}

}
