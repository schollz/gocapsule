package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/celrenheit/spider"
)

// Main data structures

type SelfData struct {
	Html string
	Date string
}

// Database functions

func (p *SelfData) Get(bucketName string) error {
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}
		k := []byte("SelfData")
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
		fmt.Printf("Could not get SelfData: %s", err)
		return err
	}
	return nil
}

func (p *SelfData) save() error {
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
			return fmt.Errorf("could not encode SelfData: %s", err)
		}
		err = bucket.Put([]byte("SelfData"), enc)
		return err
	})
	return err
}

func (p *SelfData) encode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (p *SelfData) decode(data []byte) error {
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	return nil
}

// Define the spider

var SelfSpider spider.Spider

func init() {
	SelfSpider = spider.Get("http://localhost:8080/", func(ctx *spider.Context) error {

		fmt.Println(time.Now())
		if _, err := ctx.DoRequest(); err != nil {
			return err
		}

		htmlparser, err := ctx.HTMLParser()
		if err != nil {
			return err
		}
		html, _ := htmlparser.Html()
		p := SelfData{html, time.Now().String()}

		Open()
		err = p.save()
		if err != nil {
			return fmt.Errorf("error saving SelfData")
		}
		Close()

		return nil
	})
}
