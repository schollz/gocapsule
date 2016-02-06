package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/boltdb/bolt"
	"github.com/celrenheit/spider"
)

// Main data structures

type InTheatersGroup struct {
	Data []Movie
}

type Movie struct {
	Title string
	Date  string
	Score int
	Url   string
}

// Database functions

func (p *InTheatersGroup) Get(bucketName string) error {
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}
		k := []byte("InTheatersGroup")
		val := b.Get(k)
		if val == nil {
			return nil
		}
		err = p.decode(val)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Could not get InTheatersGroup: %s", err)
		return err
	}
	return nil
}

func (p *InTheatersGroup) save() error {
	bucketName := Today()
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		enc, err := p.encode()
		if err != nil {
			return fmt.Errorf("could not encode InTheatersGroup: %s", err)
		}
		err = bucket.Put([]byte("InTheatersGroup"), enc)
		return err
	})
	return err
}

func (p *InTheatersGroup) encode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (p *InTheatersGroup) decode(data []byte) error {
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	return nil
}

// Define the spider

var InTheatersSpider spider.Spider

func init() {
	InTheatersSpider = spider.Get("http://www.metacritic.com/browse/movies/release-date/theaters/date", func(ctx *spider.Context) error {
		fmt.Println(time.Now())
		Open()
		defer Close()

		if _, err := ctx.DoRequest(); err != nil {
			return err
		}

		htmlparser, err := ctx.HTMLParser()
		if err != nil {
			return err
		}
		var p InTheatersGroup
		p.Data = []Movie{}
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
				p.Data = append(p.Data, Movie{title, date, score, url})
			}
		})
		err = p.save()
		if err != nil {
			return fmt.Errorf("error saving InTheatersGroup")
		}
		return nil
	})
}
