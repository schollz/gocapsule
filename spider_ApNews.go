package main

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/celrenheit/spider"
)

var ApNewsSpider spider.Spider

func init() {
	ApNewsSpider = spider.Get("http://hosted.ap.org/dynamic/fronts/HOME?SITE=AP&SECTION=HOME", func(ctx *spider.Context) error {
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

		htmlparser.Find(`p[class="ap-newsbriefitem-p"]`).Each(func(i int, s *goquery.Selection) {
			title := s.Find("a").Text()
			link, _ := s.Find("a").Attr("href")
			summary := s.Find(`span[class="topheadlinebody"]`).Text()
			fmt.Printf("%d: %s (%s) - %s\n", i, title, link, summary)
		})
		return nil
	})
}
