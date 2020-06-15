package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
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
	file, err := os.Create(domain + ".csv")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Comma = ','
	headline := []string{"URL", "Title", "Keywords", "Description"}
	writer.Write(headline)
	writer.Flush()

	row := make([]string, 4)
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
		title := strings.TrimSpace(e.Text)
		fmt.Println("Title: ", title)
		row[1] = title
	})

	c.OnHTML(`meta[name=description]`, func(e *colly.HTMLElement) {
		description := strings.TrimSpace(e.Attr("content"))
		fmt.Println("Description: ", description)
		row[3] = description
	})

	c.OnHTML(`meta[name=keywords]`, func(e *colly.HTMLElement) {
		keywords := strings.TrimSpace(e.Attr("content"))
		fmt.Println("Keywords: ", keywords)
		row[2] = keywords
	})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	// Before making a request
	c.OnRequest(func(r *colly.Request) {
		writer.Write(row)
		writer.Flush()
		fmt.Println("-----\nVisiting: ", r.URL)
		row[0] = r.URL.String()
	})

	defer writer.Flush()
	c.Visit("https://" + domain)
}
