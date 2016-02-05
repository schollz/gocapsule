package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/celrenheit/spider"
)

var BestMoviesInTheatersSpider spider.Spider

func init() {
	BestMoviesInTheatersSpider = spider.Get("http://www.metacritic.com/browse/movies/release-date/theaters/date", func(ctx *spider.Context) error {
		fmt.Println(time.Now())
		// Execute the request
		if _, err := ctx.DoRequest(); err != nil {
			return err
		}

		// Get goquery's html parser
		htmlparser, err := ctx.HTMLParser()
		if err != nil {
			return err
		}

		htmlparser.Find(`div[class="product_wrap group_product_wrap"]`).Each(func(i int, s *goquery.Selection) {
			title := strings.TrimSpace(s.Find("a").Text())
			score, _ := strconv.Atoi(s.Find(`div[class*="metascore_w"]`).Text())
			date := strings.TrimSpace(s.Find(`span[class="data"]`).Text())
			_, month, _ := time.Now().Date()
			curMonth := month.String()[0:3]
			_, month, _ = time.Now().Add(time.Hour * -24 * 30).Date()
			lastMonth := month.String()[0:3]
			if (strings.Contains(date, curMonth) || strings.Contains(date, lastMonth)) && score > 60 {
				urlComponent := strings.Join(strings.Split("showtimes "+title, " "), "+")
				url := "https://www.google.com/?q=" + urlComponent + "#safe=active&q=" + urlComponent
				fmt.Printf("%d: %s %d %s %s\n", i, title, score, date, url)
			}
		})
		return nil
	})
}
