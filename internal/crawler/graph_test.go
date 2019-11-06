package crawler

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewSiteGraph(t *testing.T) {
	tests := []struct {
		name string
		want SiteGraph
	}{
		{
			name: "test constructor",
			want: SiteGraph{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSiteGraph()
			if reflect.TypeOf(got) != reflect.TypeOf(SiteGraph{}) {
				t.Errorf("NewSiteGraph() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSiteGraph_AddPage(t *testing.T) {
	type fields struct {
		Graph map[string]Page
		mu    sync.RWMutex
	}
	type args struct {
		url  string
		page Page
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantLen int
	}{
		{
			name: "test1",
			fields: fields{
				Graph: map[string]Page{
					"www.website.com/page1": {
						Title: "page1",
						Links: []string{"/page2", "/page3"},
					},
				},
				mu: sync.RWMutex{},
			},
			args: args{
				url: "www.website.com/page2",
				page: Page{
					Title: "www.website.com/page2",
					Links: nil,
				},
			},
			wantLen: 2,
		},
		{
			name: "already exists",
			fields: fields{
				Graph: map[string]Page{
					"www.website.com/page1": {
						Title: "page1",
						Links: []string{"/page2", "/page3"},
					},
				},
				mu: sync.RWMutex{},
			},
			args: args{
				url: "www.website.com/page1",
				page: Page{
					Title: "page1",
					Links: []string{"/page2", "/page3"},
				},
			},
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SiteGraph{
				Graph: tt.fields.Graph,
				mu:    tt.fields.mu,
			}
			s.AddPage(tt.args.url, tt.args.page)
			pg := s.Graph[tt.args.url]
			if pg.Title != tt.args.page.Title || !reflect.DeepEqual(pg.Links, tt.args.page.Links) {
				t.Errorf("Add Page. retrieved page %v does not match inserted page %v", pg, tt.args.page)
			}
			if len(s.Graph) != tt.wantLen {
				t.Errorf("length of graph should be %d, but got %v", tt.wantLen, s.Graph)
			}
		})
	}
}

func TestSiteGraph_VisitedPage(t *testing.T) {
	type fields struct {
		Graph map[string]Page
		mu    sync.RWMutex
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Have not visited page before",
			fields:fields{
				Graph: map[string]Page{},
				mu:    sync.RWMutex{},
			},
			args: args{
				url: "www.havenotseenbefore.com",
			},
			want: false,
		},
		{
			name: "Have visited page before",
			fields:fields{
				Graph: map[string]Page{
					"www.haveseen.com": {
						Title: "something",
						Links: nil,
					},
				},
				mu:    sync.RWMutex{},
			},
			args: args{
				url: "www.haveseen.com",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SiteGraph{
				Graph: tt.fields.Graph,
				mu:    tt.fields.mu,
			}
			if got := s.VisitedPage(tt.args.url); got != tt.want {
				t.Errorf("VisitedPage() = %v, want %v", got, tt.want)
			}
		})
	}
}