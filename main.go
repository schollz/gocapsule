package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/celrenheit/spider"
	"github.com/celrenheit/spider/schedule"
	"github.com/gin-gonic/gin"
)

func main() {

	scheduler := spider.NewScheduler()
	scheduler.Add(schedule.Every(10*time.Second), ScienceNewsSpider)
	scheduler.Add(schedule.Every(10*time.Second), InTheatersSpider)
	scheduler.Add(schedule.Every(10*time.Second), GroceriesSpider)
	// scheduler.Add(schedule.Every(6*time.Second), SelfSpider)
	scheduler.Add(schedule.Every(10*time.Second), WeatherSpider)
	scheduler.Start()

	Open()
	defer Close()

	// Testing
	// Open()
	// var p WeatherData
	// p.Get(Today())
	// fmt.Println(p)
	// Close()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		day := Today()
		fmt.Println("-------")
		fmt.Println(c.ClientIP())
		fmt.Println("-------")
		c.Redirect(http.StatusMovedPermanently, c.ClientIP()+"/"+day)
	})
	router.GET("/:year/:month/:day", func(c *gin.Context) {
		year := c.Param("year")
		month := c.Param("month")
		daynum := c.Param("day")
		day := year + "/" + month + "/" + daynum
		// Open()
		var scienceNews ScienceNewsGroup
		scienceNews.Get(day)
		var inTheaters InTheatersGroup
		inTheaters.Get(day)
		var groceryList GroceryListGroup
		groceryList.Get(day)
		var weather WeatherData
		weather.Get(day)
		// Close()

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Date":        day,
			"ScienceNews": scienceNews.Data,
			"InTheaters":  inTheaters.Data,
			"Groceries":   groceryList.Data,
			"Weather":     weather,
		})
	})

	router.Run(":8001")

	// Exit 5 seconds later to let time for the request to be done.
	<-time.After(3 * time.Second)
	fmt.Println("HI")
}
