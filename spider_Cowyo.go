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

type CowyoGroup struct {
	Data []CowyoItem
}

type CowyoItem struct {
	Name string
}

// Database functions

func (p *CowyoGroup) Get(bucketName string, typeData string) error {
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}
		k := []byte(typeData)
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
		fmt.Printf("Could not get CowyoGroup: %s", err)
		return err
	}
	return nil
}

func (p *CowyoGroup) save(typeData string) error {
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
			return fmt.Errorf("could not encode CowyoGroup: %s", err)
		}
		err = bucket.Put([]byte(typeData), enc)
		return err
	})
	return err
}

func (p *CowyoGroup) encode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (p *CowyoGroup) decode(data []byte) error {
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	return nil
}

// Define the spider

var GroceriesSpider spider.Spider
var CoolListSpider spider.Spider
var TodoListSpider spider.Spider

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
		var p CowyoGroup
		p.Data = []CowyoItem{}
		dats := htmlparser.Find(`textarea `).Text()
		for _, dat := range strings.Split(dats, "\n") {
			if len(dat) > 0 {
				p.Data = append(p.Data, CowyoItem{dat})
			}
		}
		err = p.save("GroceryList")
		if err != nil {
			return fmt.Errorf("error saving CowyoGroup")
		}
		return nil
	})
	CoolListSpider = spider.Get("http://cowyo.com/cool", func(ctx *spider.Context) error {
		fmt.Print(time.Now())
		fmt.Println("CoolListSpider")

		if _, err := ctx.DoRequest(); err != nil {
			return err
		}

		htmlparser, err := ctx.HTMLParser()
		if err != nil {
			return err
		}
		var p CowyoGroup
		p.Data = []CowyoItem{}
		dats := htmlparser.Find(`textarea `).Text()
		for _, dat := range strings.Split(dats, "\n") {
			if len(dat) > 0 {
				p.Data = append(p.Data, CowyoItem{dat})
			}
		}
		err = p.save("CoolList")
		if err != nil {
			return fmt.Errorf("error saving CowyoGroup")
		}
		return nil
	})
	TodoListSpider = spider.Get("http://cowyo.com/todo", func(ctx *spider.Context) error {
		fmt.Print(time.Now())
		fmt.Println("TodoListSpider")

		if _, err := ctx.DoRequest(); err != nil {
			return err
		}

		htmlparser, err := ctx.HTMLParser()
		if err != nil {
			return err
		}
		var p CowyoGroup
		p.Data = []CowyoItem{}
		dats := htmlparser.Find(`textarea `).Text()
		for _, dat := range strings.Split(dats, "\n") {
			if len(dat) > 0 {
				p.Data = append(p.Data, CowyoItem{dat})
			}
		}
		err = p.save("TodoList")
		if err != nil {
			return fmt.Errorf("error saving TodoList")
		}
		return nil
	})
}
