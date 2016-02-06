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

type WeatherData struct {
	High int
	Low  int
	Rain int
}

// Database functions

func (p *WeatherData) Get(bucketName string) error {
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}
		k := []byte("WeatherData")
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
		fmt.Printf("Could not get WeatherData: %s", err)
		return err
	}
	return nil
}

func (p *WeatherData) save() error {
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
			return fmt.Errorf("could not encode WeatherData: %s", err)
		}
		err = bucket.Put([]byte("WeatherData"), enc)
		return err
	})
	return err
}

func (p *WeatherData) encode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (p *WeatherData) decode(data []byte) error {
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	return nil
}

// Define the spider

var WeatherSpider spider.Spider

func init() {
	WeatherSpider = spider.Get("http://forecast.weather.gov/MapClick.php?lat=35.9804&lon=-78.915&lg=english&&FcstType=digital", func(ctx *spider.Context) error {
		fmt.Println(time.Now())
		if _, err := ctx.DoRequest(); err != nil {
			return err
		}

		htmlparser, err := ctx.HTMLParser()
		if err != nil {
			return err
		}
		var p WeatherData
		p.High = 0
		p.Low = 100
		p.Rain = 0
		gotTemperature := false
		gotRain := false
		isTemperature := false
		isRain := false
		htmlparser.Find(`tr[align="center"]`).Each(func(i int, s *goquery.Selection) {
			// fmt.Println(s.Text())
			s.Find(`td`).Each(func(j int, t *goquery.Selection) {
				if strings.TrimSpace(t.Text()) == "Thunder" || strings.TrimSpace(t.Text()) == "Dewpoint (°F)" {
					isRain = false
					isTemperature = false
				}
				if isTemperature == true {
					temp, _ := strconv.Atoi(t.Text())
					if temp > p.High {
						p.High = temp
					}
					if temp < p.Low {
						p.Low = temp
					}
				}
				if strings.TrimSpace(t.Text()) == "Temperature (°F)" && gotTemperature == false {
					// fmt.Println(t.Text())
					isTemperature = true
					gotTemperature = true
				}

				if isRain == true {
					if strings.TrimSpace(t.Text()) != "--" {
						rain, _ := strconv.Atoi(t.Text())
						p.Rain += rain
					}
				}
				if strings.TrimSpace(t.Text()) == "Rain" && gotRain == false {
					// fmt.Println(t.Text())
					isRain = true
					gotRain = true
				}
			})

			// p.Data = append(p.Data, ScienceNews{title, summary, url})
		})
		Open()
		err = p.save()
		if err != nil {
			return fmt.Errorf("error saving WeatherData")
		}
		Close()
		return nil
	})
}
