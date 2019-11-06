package crawler

import (
	"golang.org/x/net/html"
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_isSameDomain(t *testing.T) {
	type args struct {
		link   string
		domain string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "same domain",
			args: args{
				link: "https://www.monzo.com/about",
				domain: "https://www.monzo.com",
			},
			want: true,
			wantErr: false,
		},
		{
			name: "same domain, without explicit host",
			args: args{
				link: "/about",
				domain: "https://www.monzo.com",
			},
			want: true,
			wantErr: false,
		},
		{
			name: "bad link, expect err",
			args: args{
				link: "12345/?#$%@#",
				domain: "https://www.monzo.com",
			},
			want: false,
			wantErr: true,
		},
		{
			name: "bad domain, expect err",
			args: args{
				link: "www.monzo.com",
				domain: "12345/?#$%@#",
			},
			want: false,
			wantErr: true,
		},
		{
			name: "different domain",
			args: args{
				link: "https://www.facebook.com",
				domain: "https://www.monzo.com",
			},
			want: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isSameDomain(tt.args.link, tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("isSameDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isSameDomain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTitle(t *testing.T) {
	type args struct {
		n *html.Node
	}

	htmlWithTitle := `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>About</title>
  </head>
  <body>
    <h1>No links here</h1>
  </body>
</html>`

	htmlWithoutTitle := `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
  </head>
  <body>
    <h1>No links here</h1>
  </body>
</html>`

	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "with title",
			html: htmlWithTitle,
			want: "About",
		},
		{
			name: "without title",
			html: htmlWithoutTitle,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.html)
			n, _ := html.Parse(reader)
			if got := getTitle(n); got != tt.want {
				t.Errorf("getTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLinks(t *testing.T) {
	doc := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Test Website</title>
</head>
<body>
  <a href="/blog.html">Blog </a>
  <a href="/sitemap.html">Sitemap</a>
  <a href="/account.html">Account</a>
  <a href="/about.html">About</a>
  <a href="https://reddit.com">Reddit External</a>
</body>
</html>`

	n, _ := html.Parse(strings.NewReader(doc))
	type args struct {
		n      *html.Node
		domain string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test1",
			args: args{
				n: n,
				domain: "www.mydomain.com",
			},
			want: []string{
				"/blog.html",
				"/sitemap.html",
				"/account.html",
				"/about.html",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLinks(tt.args.n, tt.args.domain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLinks(t *testing.T) {
	doc := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Test Website</title>
</head>
<body>
  <a href="/blog.html">Blog </a>
  <a href="/sitemap.html">Sitemap</a>
  <a href="/account.html">Account</a>
  <a href="/about.html">About</a>
  <a href="https://reddit.com">Reddit External</a>
</body>
</html>`
	type args struct {
		page   io.Reader
		domain string
	}
	tests := []struct {
		name      string
		args      args
		wantTitle string
		wantLinks []string
		wantErr   bool
	}{
		{
			name: "1",
			args: args{
				page: strings.NewReader(doc),
				domain: "www.mydomain.com",
			},
			wantTitle: "Test Website",
			wantLinks: []string{"/blog.html", "/sitemap.html", "/account.html", "/about.html"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTitle, gotLinks, err := ParseLinks(tt.args.page, tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTitle != tt.wantTitle {
				t.Errorf("ParseLinks() gotTitle = %v, want %v", gotTitle, tt.wantTitle)
			}
			if !reflect.DeepEqual(gotLinks, tt.wantLinks) {
				t.Errorf("ParseLinks() gotLinks = %v, want %v", gotLinks, tt.wantLinks)
			}
		})
	}
}