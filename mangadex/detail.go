package mangadex

import (
	"fmt"
	"mangapdf/common"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//MangaDetail ....
type MangaDetail struct {
	Title       string
	Description string
	Cover       string
	Chapters    []*Chapter
	NextHref    string
	PrivHref    string
}

//Chapter ...
type Chapter struct {
	ID   int
	Name string
	Href string
}

func (m MangaDetail) String() string {
	split := "\n=================\n"
	split2 := "\n==============\nchapters:\n"
	v := m.Title + split + m.Cover + split + m.Description + split2
	for _, v2 := range m.Chapters {
		v += split + fmt.Sprintf("%v.  ", v2.ID) + v2.Name
	}
	return v
}

//Detail ...
func (m *Mangadex) Detail(href string) (*MangaDetail, error) {
	resp, err := m.fetch(req{
		link:   href,
		method: "GET",
	})
	doc, err := common.GetDoc(resp)
	if err != nil {
		return nil, err
	}
	md := &MangaDetail{}
	md.Cover, _ = doc.Find("img.rounded").Attr("src")
	md.Title = doc.Find("h6 > span.mx-1").Text()
	l1 := doc.Find("div.col-lg-3.col-xl-2.strong")
	for i := 0; i < l1.Length(); i++ {
		v := l1.Eq(i).Text()
		if v == "Description:" {
			d := l1.Eq(i).Next().Text()
			md.Description = strings.Split(d, "\n")[0]
		}
	}
	l2 := doc.Find(".chapter-container > .row.no-gutters")
	md.Chapters = make([]*Chapter, 0)
	id := 0
	l2.Each(func(i int, s *goquery.Selection) {
		ls := s.Find(".rounded.flag.flag-gb").Length()
		if ls == 0 {
			return
		}
		c := &Chapter{}
		c.ID = id
		c.Name = s.Find("a.text-truncate").Text()
		c.Href, _ = s.Find("a.text-truncate").Attr("href")
		md.Chapters = append(md.Chapters, c)
		id++
	})
	return md, nil
}
