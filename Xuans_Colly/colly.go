package Xuans_Colly

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
	"github.com/velebak/colly-sqlite3-storage/colly/sqlite3"
)

func XuansColly() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.example.com"),
	)
	storage := &sqlite3.Storage{
		Filename: "./qsy.db",
	}
	defer storage.Close()

	err := c.SetStorage(storage)
	if err != nil {
		panic(err)
	}

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	q, _ := queue.New(8, storage)
	q.AddURL("https://www.example.com")

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		q.AddURL(e.Request.AbsoluteURL(e.Attr("href")))
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.Request.URL, "\t", r.StatusCode)
	})

	q.Run(c)
	log.Println(c)
}
