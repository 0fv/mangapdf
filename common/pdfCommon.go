package common

import (
	"bytes"
	"image"
	"log"

	"github.com/signintech/gopdf"
)

//AddPictureBytesToPdf ....
func AddPictureBytesToPdf(b []byte, pdf *gopdf.GoPdf) error {
	pdf.AddPage()
	if len(b) == 0 {
		log.Println("empty page")
		return nil
	}
	img, _, err := image.Decode(bytes.NewBuffer(b))
	if err != nil {
		return err
	}
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
	pic, err := gopdf.ImageHolderByBytes(b)
	if err != nil {
		return err
	}
	pdf.ImageByHolder(pic, ws, hs, p)
	return nil
}
