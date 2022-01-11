# Scraper for Joel Spolsky's Blog https://www.joelonsoftware.com/
Joel Spolsky is the former CEO of Stack Overflow and has his own blog. He doesn't have a list of his posts with just the title and links to the post; this archive page https://www.joelonsoftware.com/archives/ shows the months (not titles) he's posted in, and clicking on the months shows all the actual article content for that month (instead of only the titles). This scraper extracts those article titles + links and outputs them to an HTML file for easier browsing. It's also written in different languages to act as practice and code samples.

# Golang
Original version.

## Running
* `cd blog-scrape/golang/`
* `go run main.go scrape` to create local text file of posts
* `go run main.go template` to generate html
* Open `dist/generated.html` to see list of articles

## Config
* Open `config.json`
* `MaxGoRoutines` sets the maximum number of goroutines used to scrape the blog
* `BufferSize` sets the capacity of the channel and array used to store articles
* `ScrapeDelay` sets the delay in milliseconds before a goroutine scrapes another link

# Python
Written using Python 3.9. Doesn't have HTML file feature (yet?).

## Running
* `cd blog-scrape/python/`
* `python scrape.py`
* See output in `dist/scraped-links.json`

## Config
* See `config.py` for options and documentation
