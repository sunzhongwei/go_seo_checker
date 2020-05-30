package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"os"
	"strings"
	"time"
)

func main() {
	var domain string
	if len(os.Args) != 2 {
		fmt.Println("Please input domain! e.g. ./go_seo_checker www.notefeel.com")
		return
	}
	domain = os.Args[1]

	c := colly.NewCollector(
		colly.AllowedDomains(domain),
	)
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	c.Limit(&colly.LimitRule{
		DomainGlob:  domain,
		Parallelism: 1,
		Delay:       2 * time.Second,
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("Title: ", strings.TrimSpace(e.Text))
	})

	c.OnHTML(`meta[name=description]`, func(e *colly.HTMLElement) {
		fmt.Println("Description: ", strings.TrimSpace(e.Attr("content")))
	})

	c.OnHTML(`meta[name=keywords]`, func(e *colly.HTMLElement) {
		fmt.Println("Keywords: ", strings.TrimSpace(e.Attr("content")))
	})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	// Before making a request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("-----\nVisiting: ", r.URL)
	})

	c.Visit("https://" + domain)
}
