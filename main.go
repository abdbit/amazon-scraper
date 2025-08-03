package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

// custom type
type Phone struct {
	Title string `csv:"title"`
	Price string `csv:"price"`
	Link  string `csv:"link"`
}

var phones []Phone

func main() {
	appendToFile("./iphones.csv", "Title~Price~Link\n") // adding header row

	c := colly.NewCollector()

	c.OnHTML("div[role=listitem]", ExtractPhoneDetails)

	pageNumber := 0

	for pageNumber < 20 {
		pageNumber = pageNumber + 1
		url := fmt.Sprintf("https://www.amazon.com/s?k=iPhone+15&page=%d", pageNumber)
		c.Visit(url)
	}

	fmt.Printf("\n\nTotal number of products: %d\n\n", len(phones))
}

func ExtractPhoneDetails(h *colly.HTMLElement) {
	phone := Phone{}

	phone.Title = strings.TrimSpace(h.DOM.Find("div[data-cy=title-recipe] h2.a-size-medium").Text())
	if phone.Title == "" {
		return
	}

	phone.Price = strings.TrimSpace(h.DOM.Next().Find("span[class=a-price] span[class=a-offscreen]").Text())
	phone.Link = "https://www.amazon.com" + strings.TrimSpace(h.DOM.Find("a.a-link-normal").AttrOr("href", ""))

	phones = append(phones, phone)

	fmt.Printf("---> %s~%s~%s\n", phone.Title, phone.Price, phone.Link)

	csvRow := fmt.Sprintf("%s~%s~%s\n", phone.Title, phone.Price, phone.Link)
	appendToFile("./iphones.csv", csvRow)
}

func appendToFile(filename string, data string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
		panic(err)
	}
}
