package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/boltdb/bolt"
	"github.com/celrenheit/spider"
)

// Main data structures

type ScienceNewsGroup struct {
	Data []ScienceNews
}

type ScienceNews struct {
	Title   string
	Summary string
	Url     string
}

// Database functions

func (p *ScienceNewsGroup) Get(bucketName string) error {
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}
		k := []byte("ScienceNewsGroup")
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
	fmt.Println(p)
	if err != nil {
		fmt.Printf("Could not get ScienceNewsGroup: %s", err)
		return err
	}
	return nil
}

func (p *ScienceNewsGroup) save() error {
	bucketName := time.Now().String()[0:10]
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
			return fmt.Errorf("could not encode ScienceNewsGroup: %s", err)
		}
		err = bucket.Put([]byte("ScienceNewsGroup"), enc)
		return err
	})
	return err
}

func (p *ScienceNewsGroup) encode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (p *ScienceNewsGroup) decode(data []byte) error {
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	return nil
}

// Define the spider

var ScienceNewsSpider spider.Spider

func init() {
	ScienceNewsSpider = spider.Get("http://www.sciencemag.org/news", func(ctx *spider.Context) error {
		fmt.Println(time.Now())
		if _, err := ctx.DoRequest(); err != nil {
			return err
		}

		htmlparser, err := ctx.HTMLParser()
		if err != nil {
			return err
		}
		var p ScienceNewsGroup
		p.Data = []ScienceNews{}
		htmlparser.Find(`div[class="media__body"]`).Each(func(i int, s *goquery.Selection) {
			title := strings.TrimSpace(s.Find("h2").Text())
			url, _ := s.Find("h2 > a").Attr("href")
			url = "http://www.sciencemag.org" + url
			summary := "None"
			fmt.Printf("%d: %s (%s) - %s\n", i, title, url, summary)
			p.Data = append(p.Data, ScienceNews{title, summary, url})
		})
		err = p.save()
		if err != nil {
			return fmt.Errorf("error saving ScienceNewsGroup")
		}
		return nil
	})
}
