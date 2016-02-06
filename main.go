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
	// scheduler.Add(schedule.Every(10*time.Second), ScienceNewsSpider)
	// scheduler.Add(schedule.Every(10*time.Second), InTheatersSpider)
	// scheduler.Add(schedule.Every(2*time.Second), GroceriesSpider)
	scheduler.Add(schedule.Every(6*time.Second), SelfSpider)
	scheduler.Start()

	// Testing
	// Open()
	// defer Close()
	// var p GroceryListGroup
	// p.Get(Today())
	// fmt.Println(p)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		day := Today()

		Open()
		var scienceNews ScienceNewsGroup
		scienceNews.Get(day)
		var inTheaters InTheatersGroup
		inTheaters.Get(day)
		var groceryList GroceryListGroup
		groceryList.Get(day)
		Close()
		fmt.Println(groceryList)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Date":        day,
			"ScienceNews": scienceNews.Data,
			"InTheaters":  inTheaters.Data,
			"Groceries":   groceryList.Data,
		})
	})
	router.Run(":8080")

	// Exit 5 seconds later to let time for the request to be done.
	<-time.After(3 * time.Second)
	fmt.Println("HI")
}
