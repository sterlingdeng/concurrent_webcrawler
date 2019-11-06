package crawler

import (
	"golang.org/x/net/html"
	"io"
	"net/url"
)

func ParseLinks(page io.Reader, domain string) (title string, links []string, err error) {
	node, err := html.Parse(page)
	if err != nil {
		return "", nil, err
	}
	return getTitle(node), getLinks(node, domain), nil
}

func getLinks(n *html.Node, domain string) []string {
	var links []string
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				same, err := isSameDomain(attr.Val, domain)
				if err != nil {
					continue
				}
				if same {
					return []string{attr.Val}
				}
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		l := getLinks(child, domain)
		links = append(links, l...)
	}
	return links
}

func isSameDomain(link, domain string) (bool, error) {
	lurl, err := url.Parse(link)
	if err != nil {
		return false, err
	}
	durl, err := url.Parse(domain)
	if err != nil {
		return false, err
	}
	if !lurl.IsAbs() || lurl.Hostname() == durl.Hostname() {
		return true, nil
	}
	return false, nil
}

func getTitle(n *html.Node) string {
	var title string
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		title = getTitle(child)
		if title != "" {
			break
		}
	}
	return title
}
