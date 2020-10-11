package mangadex

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"testing"

	"github.com/signintech/gopdf"
)

func TestLogin(t *testing.T) {
	m := Mangadex{}
	err := m.Login("xxxx", "xxxx")
	log.Println(err)
}

func TestSearch(t *testing.T) {
	m := Mangadex{}
	log.Println(m.checkFile())
	l, err := m.Search("one", 1)
	log.Println(err)
	for _, v := range l {
		fmt.Println(v)
	}
}

func TestDetail(t *testing.T) {
	m := Mangadex{}
	log.Println(m.checkFile())
	l, err := m.Detail("/title/21562/kusuriya-no-hitorigoto")
	log.Println(err)
	fmt.Println(l)
}

func TestChapter(t *testing.T) {
	m := Mangadex{}
	log.Println(m.checkFile())
	c, _ := m.GetContent("24408")
	endChan := make(chan int, 1)
	m.ToPDF(c, endChan)
	select {
	case <-endChan:
		return
	}
}

func TestPdf(t *testing.T) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	file, _ := ioutil.ReadFile("book.jpg")

	img, _, _ := image.Decode(bytes.NewBuffer(file))
	x := img.Bounds().Dx()
	y := img.Bounds().Dy()
	var p *gopdf.Rect
	var ws float64 = 0
	var hs float64 = 0
	if x > y {
		w := float64(595)
		h := float64(y) / float64(x) * float64(w)
		p = &gopdf.Rect{W: w, H: h}
		hs = (842 - h) / 2
	} else {
		h := float64(842)
		w := float64(x) / float64(y) * float64(h)
		p = &gopdf.Rect{W: w, H: h}
		ws = (595 - w) / 2
	}
	// var PageSizeA4 = &Rect{W: 595, H: 842, unitOverride: UnitPT}
	pic, _ := gopdf.ImageHolderByBytes(file)
	pdf.ImageByHolder(pic, ws, hs, p)
	// pdf.SetX(250) //move current location
	// pdf.SetY(200)

	pdf.WritePdf("image.pdf")
}
