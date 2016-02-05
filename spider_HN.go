package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/celrenheit/spider"
)

var HackerNewsSpider spider.Spider

func init() {
	HackerNewsSpider = spider.Get("https://news.ycombinator.com/", func(ctx *spider.Context) error {
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

		title := ""
		origin := ""
		url := ""
		score := 0
		haveTitle := false
		htmlparser.Find(`tr`).Each(func(i int, s *goquery.Selection) {
			// fmt.Println(i)
			if math.Mod(float64(i), 2) == 0 {
				title = s.Find(`td[class="title"] > a`).Text()
				origin = s.Find(`span[class="sitebit comhead"]`).Text()
				url, _ = s.Find(`td[class="title"] > a`).Attr("href")
				haveTitle = true
			} else if math.Mod(float64(i), 2) == 1 && haveTitle == true {
				score, _ = strconv.Atoi(strings.Split(s.Find(`span[class="score"]`).Text(), " ")[0])
				if score > 10 {
					fmt.Println(title, origin, url, score)
				}
				haveTitle = false
			}
		})
		return nil
	})
}
