package main

import (
	"flag"
	"log"
	"mangapdf/mangadex"
)

var (
	chapter = flag.String("c", "", "chapterID")
)

func main() {
	flag.Parse()
	m := mangadex.Mangadex{}
	if chapter != nil && *chapter != "" {
		c, err := m.GetContent(*chapter)
		if err != nil {
			log.Println(err)
			return
		}
		endChan := make(chan int, 1)
		m.ToPDF(c, endChan)
		select {
		case <-endChan:
			return
		}
	}

}
