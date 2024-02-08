package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Book struct {
	Title string
	Price string
}

func main() {
	file, err := os.Create("export.csv")
	if err != nil {
		log.Fatal(err)

	}
	defer file.Close()

	//creating csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()
	//writing the header
	headers := []string{"Title", "Price"}
	writer.Write(headers)

	c := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"),
	)

	proxyUsername := "<Username>"
	proxyPassword := "<Password>"
	proxyUrl := fmt.Sprintf("http://customer-%s:%s@pr.oxylabs.io:7777", proxyUsername, proxyPassword)

	c.SetProxy(proxyUrl)

	c.OnHTML(".product_pod", func(e *colly.HTMLElement) {
		book := Book{}
		book.Title = e.ChildAttr(".image_container img", "alt")
		book.Price = e.ChildText(".price_color")
		row := []string{book.Title, book.Price}
		writer.Write(row)
	})

	c.OnHTML(".next a", func(h *colly.HTMLElement) {
		nextPage := h.Request.AbsoluteURL(h.Attr("href"))
		c.Visit(nextPage)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://books.toscrape.com/")

}
