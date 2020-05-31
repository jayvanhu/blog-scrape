package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"scrape/util"
	str "strings"
	"sync"
	"time"
)

/// Globals
var wg sync.WaitGroup
var scrapeFile string = "dist/scraped-links.txt"
var client *http.Client

type Config struct {
	MaxGoroutines int
	BufferSize    int
	ScrapeDelay   time.Duration
}

var config Config

func loadConfig() {
	file, err := os.Open("config.json")
	util.HandleErr(err)
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	util.HandleErr(err)

	json.Unmarshal(bytes, &config)

	config.ScrapeDelay = config.ScrapeDelay * time.Millisecond
}

/// Main
type Article struct {
	Date  string
	Title string
	Href  string
}

func newReq(uri string) *http.Request {
	req, err := http.NewRequest("GET", uri, nil)
	util.HandleErr(err)

	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent/Firefox
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:10.0) Gecko/20100101 Firefox/10.0")
	return req
}

func getSiteLinks() []string {
	site := "https://www.joelonsoftware.com/archives/"

	req := newReq(site)

	res, err := client.Do(req)
	util.HandleErr(err)
	defer res.Body.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	util.HandleErr(err)

	links := make([]string, 0, config.BufferSize)
	document.Find(".yearly-archive > .month > h3 > a").Each(func(i int, el *goquery.Selection) {
		href, ok := el.Attr("href")
		if ok {
			links = append(links, href)
		}
	})
	return links
}

func scrapeLinks(links []string, output chan Article, maxGoRoutines int) {
	hrefs := make(chan string, config.BufferSize)

	fmt.Println("Creating goroutines")
	for i := 0; i < maxGoRoutines; i++ {
		wg.Add(1)
		go scrapeLink(hrefs, output, true)
	}

	fmt.Println("Sending links")
	for _, link := range links {
		hrefs <- link
	}

	fmt.Println("Closing hrefs chan")
	close(hrefs)
}

func scrapeLink(hrefs <-chan string, output chan<- Article, wait bool) {
	for {
		fmt.Println("Scraping link")
		link, more := <-hrefs
		if !more {
			fmt.Println("Hrefs is closed")
			wg.Done()
			return
		}

		req := newReq(link)
		res, err := client.Do(req)
		util.HandleErr(err)
		defer res.Body.Close()
		document, err := goquery.NewDocumentFromReader(res.Body)

		document.Find("header.entry-header").Each(func(i int, header *goquery.Selection) {
			dateStr, okDate := header.Find("time.entry-date").Attr("datetime")

			anchor := header.Find("h1.entry-title > a")
			href, okHref := anchor.Attr("href")
			title := anchor.Text()

			if okDate && okHref {
				article := Article{
					Href:  href,
					Title: title,
					Date:  dateStr,
				}
				output <- article
			} else {
				fmt.Println("scrapeLink() :: field missing: ", dateStr, title, href)
			}
		})
		if wait {
			time.Sleep(config.ScrapeDelay)
		}
	}
}

func scrape() {
	client = &http.Client{
		Timeout: 10 * time.Second,
	}

	links := getSiteLinks()
	// Uncomment to limit links
	// links = links[:2]

	// future: sort results by date
	initialArticleBuffer := 500
	articles := make([]Article, 0, initialArticleBuffer)
	articleCh := make(chan Article, initialArticleBuffer)

	wait := make(chan bool)
	go func() {
		for {
			fmt.Println("Reading from articleCh")
			article, more := <-articleCh
			if more {
				articles = append(articles, article)
				fmt.Println("Recv from articleCh", article)
			} else {
				fmt.Println("articleCh is closed")
				break
			}
		}
		wait <- true
	}()
	// future: not obvious this uses goroutine
	scrapeLinks(links, articleCh, config.MaxGoroutines)
	wg.Wait()
	fmt.Println("Closing articleCh")
	close(articleCh)
	<-wait

	bytes, err := json.Marshal(articles)
	util.HandleErr(err)
	util.WriteToFile(scrapeFile, string(bytes))
}

func createHtml() {
	file, err := os.Open(scrapeFile)
	util.HandleErr(err)
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	util.HandleErr(err)

	var articles []Article
	json.Unmarshal(bytes, &articles)

	render, err := template.ParseFiles("template.html")
	util.HandleErr(err)

	outputFile, err := os.Create("dist/generated.html")
	util.HandleErr(err)
	defer outputFile.Close()

	render.Execute(outputFile, articles)
	render.Execute(os.Stdout, articles)
}

func main() {
	fmt.Println("Start")
	startTime := time.Now()
	loadConfig()

	if len(os.Args) > 1 {
		cmd := str.ToLower(os.Args[1])
		switch cmd {
		case "scrape":
			scrape()
		case "template":
			createHtml()
		default:
			fmt.Println("Invalid command-line argument")
		}
	} else {
		fmt.Println("No command-line argument provided")
	}

	secs := time.Since(startTime).Seconds()
	fmt.Println("Finished in seconds:", secs)
}
