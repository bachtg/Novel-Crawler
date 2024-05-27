package internal

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"novel_crawler/internal/category"
	"strings"
	"time"
)

type Source interface {
	GetCategories() []*category.Category
}

type TruyenFull struct {
}

func (truyenFull *TruyenFull) GetCategories() []*category.Category {

	var listCategories []*category.Category

	c := colly.NewCollector(
		colly.AllowedDomains("truyenfull.vn"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.Async(true),
		colly.MaxDepth(1),
	)

	// Set up event handlers
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML(".dropdown-menu a", func(e *colly.HTMLElement) {
		// Printing all URLs associated with the a links in the page
		url := e.Attr("href")
		name := e.Text

		if strings.Contains(url, "the-loai") {
			listCategories = append(listCategories, &category.Category{
				Url:  url,
				Name: name,
			})
		}

	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, "scraped!")
	})

	// Adding a delay between requests
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	// Now visit the page
	err := c.Visit("https://truyenfull.vn/")
	if err != nil {
		fmt.Println("err->", err)
	}

	// Wait until all asynchronous tasks are finished
	c.Wait()
	return listCategories
}
