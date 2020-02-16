package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	urls := mainPage()

	subPage(urls)
}

func mainPage() (urls map[int]string) {
	fName := "scheme.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"id", "scheme"})

	// Instantiate default collector
	c := colly.NewCollector()

	urls = make(map[int]string)
	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, element *colly.HTMLElement) {
			row := element.ChildTexts("td")
			writer.Write(row)

			index, _ := strconv.Atoi(row[0])
			url := element.Request.AbsoluteURL(element.ChildAttr("a", "href"))
			urls[index] = url
		})
	})

	c.Visit("https://bintable.com/card-schemes")
	log.Printf("Scraping finished, check file %q for results\n", fName)
	return urls
}

func subPage(urls map[int]string) {
	fName := "schemeDetails.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"bin", "scheme", "type", "category", "issuer"})

	// Instantiate default collector
	c := colly.NewCollector()
	c.SetRequestTimeout(120 * time.Second)

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, element *colly.HTMLElement) {
			row := element.ChildTexts("td")
			writer.Write(row)
			fmt.Println(row)
		})
	})

	c.OnHTML(".pagination li:last-child a", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
		fmt.Println("start visiting next page: ", nextPage)
		e.Request.Visit(nextPage)
	})

	for _, v := range urls {
		c.Visit(v)
	}
	log.Printf("Scraping finished, check file %q for results\n", fName)
}