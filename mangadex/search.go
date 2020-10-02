package mangadex

import (
	"fmt"
	"mangapdf/common"

	"github.com/PuerkitoBio/goquery"
)

var search = "/search"

//SearchResult ...
type SearchResult struct {
	ID           int
	Title        string
	Introduction string
	DataID       string
	Href         string
}

func (s SearchResult) String() string {
	split := "\n============\n"
	split2 := "\n==============\n\n"
	v := s.Introduction
	if len(s.Introduction) > 300 {
		v = s.Introduction[:300] + "..."
	}
	return fmt.Sprint(s.ID) + "." + s.Title + split + v + split2
}

//Search ...
func (m *Mangadex) Search(keyword string, page int) ([]*SearchResult, error) {
	resp, err := m.fetch(req{
		link: search,
		linkParam: map[string]string{
			"title": keyword,
			"s":     "0",
			"p":     fmt.Sprint(page),
		},
	})
	if err != nil {
		return nil, err
	}
	doc, err := common.GetDoc(resp)
	if err != nil {
		return nil, err
	}
	l := doc.Find(".manga-entry")
	var r []*SearchResult = make([]*SearchResult, 0, l.Length())
	l.Each(func(i int, s *goquery.Selection) {
		sr := SearchResult{}
		sr.DataID, _ = s.Attr("data-id")
		sr.Title = s.Find(".manga_title").Text()
		sr.Href, _ = s.Find(".manga_title").Attr("href")
		sr.Introduction = s.Find(".list-inline.m-1").Next().Text()
		sr.ID = i
		r = append(r, &sr)
	})
	return r, nil
}
