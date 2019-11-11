package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
	"webcrawler/internal/crawler"
)

const (
	DefaultDomain      = "http://localhost:8080"
	DefaultWorkerCount = 10
	DefaultLogLevel    = "info"
	DefaultMaxDepth    = 3
)

var (
	DomainUrl   string
	WorkerCount int
	LogLevel    string
	MaxDepth    int
)

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawls a given website and returns the sitemap",
	Long:  "Crawls a given website, if none is provided, it will crawl localhost:8080",
	Run: func(cmd *cobra.Command, args []string) {
		sitemap := Crawl()
		for url, pageInfo := range sitemap {
			for _, link := range pageInfo.Links {
				fmt.Printf("Parent URL: %s. Child URL: %s\n", url, link)
			}
		}
	},
}

func init() {
	crawlCmd.Flags().StringVarP(&DomainUrl, "domain", "d", DefaultDomain, "Domain for webcrawler to crawl")
	crawlCmd.Flags().IntVarP(&WorkerCount, "workers", "w", DefaultWorkerCount, "Sets the number of concurrent workers")
	crawlCmd.Flags().StringVarP(&LogLevel, "log", "l", DefaultLogLevel, "Sets the logging level")
	crawlCmd.Flags().IntVarP(&MaxDepth, "depth", "n", DefaultMaxDepth, "Max depth the crawler will crawl.")
	rootCmd.AddCommand(crawlCmd)
}

func Crawl() map[string]crawler.Page {
	t1 := time.Now()
	defer func() { fmt.Println("Time:", time.Now().Sub(t1)) }()

	c := crawler.NewCrawler(DomainUrl, WorkerCount, MaxDepth, LogLevel)
	err := c.Crawl()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	return c.GetRelationships()
}
