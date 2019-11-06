package crawler

import "sync"

type Page struct {
	Title string
	Links []string
}

type SiteGraph struct {
	Graph map[string]Page
	mu    sync.RWMutex
}

func NewSiteGraph() SiteGraph {
	return SiteGraph{
		Graph: make(map[string]Page),
	}
}

func (s *SiteGraph) AddPage(url string, page Page) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exist := s.Graph[url]; !exist {
		s.Graph[url] = page
	}
}

func (s *SiteGraph) VisitedPage(url string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, visited := s.Graph[url]; visited {
		return true
	}
	return false
}

//func (s *SiteGraph) addPageData(pages ...Page) {
//	for _, p := range pages {
//		if _, exist := s.PageData[p.Url]; !exist {
//			s.PageData[p.Url] = p
//		}
//	}
//}
