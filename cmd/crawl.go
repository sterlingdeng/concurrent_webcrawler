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
)

var (
	DomainUrl   string
	WorkerCount int
	LogLevel    string
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
	crawlCmd.Flags().StringVarP(&LogLevel, "log", "l", DefaultLogLevel, "Sets the logging level")
	rootCmd.AddCommand(crawlCmd)
}

func Crawl() {
	t1 := time.Now()
	defer func() {fmt.Println("Time:", time.Now().Sub(t1))}()

	c := crawler.NewCrawler(DomainUrl, WorkerCount, LogLevel)
	err := c.Crawl()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	fmt.Println("Sitemap:", c.GetRelationships())
}
