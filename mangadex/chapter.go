package mangadex

import (
	"fmt"
	"io/ioutil"
	"log"
	"mangapdf/common"
	"net/http"
	"sort"
	"time"

	"github.com/signintech/gopdf"
)

//ChapterContent ...
type ChapterContent struct {
	ID        int      `json:"id"`
	Hash      string   `json:"hash"`
	Volume    string   `json:"volume"`
	Chapter   string   `json:"chapter"`
	Title     string   `json:"title"`
	LangName  string   `json:"lang_name"`
	Server    string   `json:"server"`
	PageArray []string `json:"page_array"`
}

func (c *ChapterContent) genLink(i int) string {
	return c.Server + c.Hash + "/" + c.PageArray[i]
}

func (c *ChapterContent) getLinkLen() int {
	return len(c.PageArray)
}

func (c ChapterContent) String() string {
	return "v." + c.Volume + "c." + c.Chapter + " " + c.Title
}

var linkCapter = "https://mangadex.org/api/?id=%v&server=null&saver=0&type=chapter"

//GetContent ...
func (m *Mangadex) GetContent(chapter string) (*ChapterContent, error) {
	// resp, err := m.fetch(req{
	// 	method: "Get",
	// 	link:   fmt.Sprintf(linkCapter, capter),
	// 	header: map[string]string{
	// 		"accept": "application/json",
	// 	},
	// })
	resp, err := http.Get(fmt.Sprintf(linkCapter, chapter))
	if err != nil {
		return nil, err
	}
	var cc ChapterContent
	err = common.GetBody(resp, &cc)
	return &cc, err
}

type picBuf struct {
	Index int
	Link  string
	Buf   []byte
}

func (p *picBuf) getPicBuf(m *Mangadex) {
	log.Println("start page", p.Index)
	start := time.Now()
	defer func() {
		end := time.Now().Sub(start).Seconds()
		log.Println("end page", p.Index, "using time", end, "s")
	}()
	var resp *http.Response
	var err error
	for i := 0; i < 3; i++ {
		resp, err = m.fetch(req{
			method: "GET",
			link:   p.Link,
		})
		// resp, err = http.Get(p.Link)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 3)
		} else {
			break
		}
	}
	if err != nil {
		log.Println("error:", err)
		return
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("io error", err)
		return
	}
	p.Buf = buf
}

type picBufs []picBuf

//Len()
func (s picBufs) Len() int {
	return len(s)
}

//Less():
func (s picBufs) Less(i, j int) bool {
	return s[i].Index < s[j].Index
}

//Swap()
func (s picBufs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//ToPDF ...
func (m *Mangadex) ToPDF(c *ChapterContent, endCh chan int) {
	ch := make(chan picBuf, 10)
	ch2 := make(chan picBuf, c.getLinkLen())
	fileName := fmt.Sprint(c) + ".pdf"
	go m.getBufList(&ch2, fileName, c.getLinkLen(), endCh)
	go m.getPicBuf(&ch, &ch2)
	for i := range c.PageArray {
		l := c.genLink(i)
		pic := picBuf{}
		pic.Index = i
		pic.Link = l
		ch <- pic
	}
}

//ToPDF to pdf
func (m *Mangadex) getBufList(ch *chan picBuf, fileName string, length int, endCh chan int) {
	var tmpBuf picBufs = make(picBufs, 0, length)
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	defer func() {
		log.Println(fileName, "complete")
		endCh <- 1
	}()
	sum := 0
	for v := range *ch {
		tmpBuf = append(tmpBuf, v)
		sum++
		if sum == length {
			sort.Sort(tmpBuf)
			for _, pic := range tmpBuf {
				err := common.AddPictureBytesToPdf(pic.Buf, &pdf)
				if err != nil {
					log.Println(err)
					return
				}
			}
			pdf.WritePdf(fileName)
			close(*ch)
			return
		}
	}
}

func (m *Mangadex) getPicBuf(ch1 *chan picBuf, ch2 *chan picBuf) {
	runingChan := make(chan int, 10)
	for v := range *ch1 {
		runingChan <- 1
		go func(v2 picBuf) {
			v2.getPicBuf(m)
			*ch2 <- v2
			<-runingChan
		}(v)
	}
}
