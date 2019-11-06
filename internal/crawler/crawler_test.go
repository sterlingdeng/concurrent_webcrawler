package crawler

import (
	"reflect"
	"testing"

	"webcrawler/internal"
	"webcrawler/internal/crawler/test"

	log "github.com/sirupsen/logrus"
)

func TestCrawler_Crawl(t *testing.T) {
	ts := test.NewServer()
	defer ts.Stop()

	type fields struct {
		domain      string
		graph       SiteGraph
		logger      *log.Logger
		workerCount int
	}
	tests := []struct {
		name      string
		fields    fields
		wantErr   bool
		wantGraph map[string]Page
	}{
		{
			name: "with mock html",
			fields: fields{
				domain:      ts.Ts.URL,
				graph:       NewSiteGraph(),
				logger:      internal.NewLogger("DEBUG"),
				workerCount: 3,
			},
			wantErr: false,
			wantGraph: map[string]Page{
				ts.Ts.URL:                       {"Test Website", []string{"/blog.html", "/sitemap.html", "/account.html", "/about.html"}},
				ts.Ts.URL + "/about.html":       {"About", nil},
				ts.Ts.URL + "/account.html":     {"Account", []string{"/checking.html", "/investments.html", "/savings.html"}},
				ts.Ts.URL + "/blog.html":        {"Blog", []string{"/stories.html"}},
				ts.Ts.URL + "/checking.html":    {"Checking Account Page", nil},
				ts.Ts.URL + "/index.html":       {"Test Website", []string{"/blog.html", "/sitemap.html", "/account.html", "/about.html"}},
				ts.Ts.URL + "/investments.html": {"Investments", nil},
				ts.Ts.URL + "/savings.html":     {"Savings", []string{"/index.html"}},
				ts.Ts.URL + "/sitemap.html":     {"Sitemap", nil},
				ts.Ts.URL + "/stories.html":     {"Stories", nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{
				domain:      tt.fields.domain,
				graph:       tt.fields.graph,
				logger:      tt.fields.logger,
				workerCount: tt.fields.workerCount,
			}
			if err := c.Crawl(); (err != nil) != tt.wantErr {
				t.Errorf("Crawl() error = %v, wantErr %v", err, tt.wantErr)
			}
			relationships := c.GetRelationships()
			if !reflect.DeepEqual(relationships, tt.wantGraph) {
				t.Errorf("Crawl() output error. want %v, got %v", tt.wantGraph, relationships)

			}
		})
	}
}

func TestNewCrawler(t *testing.T) {
	type args struct {
		domain   string
		count    int
		loglevel string
	}
	tests := []struct {
		name string
		args args
		want Crawler
	}{
		{
			name: "test constructor",
			args: args{
				domain:   "test.com",
				count:    1,
				loglevel: "DEBUG",
			},
			want: Crawler{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCrawler(tt.args.domain, tt.args.count, tt.args.loglevel)
			if reflect.TypeOf(got) != reflect.TypeOf(Crawler{}) {
				t.Errorf("NewCrawler() = %v, want %v", got, tt.want)
			}
		})
	}
}
