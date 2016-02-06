package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/celrenheit/spider"
)

// Main data structures

type GroceryListGroup struct {
	Data []GroceryItem
}

type GroceryItem struct {
	Name string
}

// Database functions

func (p *GroceryListGroup) Get(bucketName string) error {
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}
		k := []byte("GroceryListGroup")
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
		fmt.Printf("Could not get GroceryListGroup: %s", err)
		return err
	}
	return nil
}

func (p *GroceryListGroup) save() error {
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
			return fmt.Errorf("could not encode GroceryListGroup: %s", err)
		}
		err = bucket.Put([]byte("GroceryListGroup"), enc)
		return err
	})
	return err
}

func (p *GroceryListGroup) encode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (p *GroceryListGroup) decode(data []byte) error {
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	return nil
}

// Define the spider

var GroceriesSpider spider.Spider

func init() {
	GroceriesSpider = spider.Get("http://cowyo.com/grocerylist", func(ctx *spider.Context) error {
		fmt.Print(time.Now())
		fmt.Println("GroceriesSpider")

		if _, err := ctx.DoRequest(); err != nil {
			return err
		}

		htmlparser, err := ctx.HTMLParser()
		if err != nil {
			return err
		}
		var p GroceryListGroup
		p.Data = []GroceryItem{}
		dats := htmlparser.Find(`textarea `).Text()
		for _, dat := range strings.Split(dats, "\n") {
			if len(dat) > 0 {
				p.Data = append(p.Data, GroceryItem{dat})
			}
		}
		// Open()
		err = p.save()
		// Close()
		if err != nil {
			return fmt.Errorf("error saving GroceryListGroup")
		}
		return nil
	})
}
