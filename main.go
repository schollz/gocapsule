package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/celrenheit/spider"
	"github.com/celrenheit/spider/schedule"
	"github.com/gin-gonic/gin"
)

var ExternalIP string

func main() {
	if len(os.Args) == 1 {
		log.Fatal("You need to specify the external IP address")
	}
	ExternalIP = os.Args[1]
	loadSelfSpider()
	scheduler := spider.NewScheduler()
	// scheduler.Add(schedule.Every(600*time.Second), ScienceNewsSpider)
	scheduler.Add(schedule.Every(600*time.Second), InTheatersSpider)
	scheduler.Add(schedule.Every(600*time.Second), GroceriesSpider)
	scheduler.Add(schedule.Every(600*time.Second), CoolListSpider)
	scheduler.Add(schedule.Every(600*time.Second), TodoListSpider)
	scheduler.Add(schedule.Every(600*time.Second), SelfSpider)
	scheduler.Add(schedule.Every(600*time.Second), WeatherSpider)
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
		c.Redirect(307, "http://"+ExternalIP+"/"+day)
	})
	router.GET("/:year/:month/:day", func(c *gin.Context) {
		year := c.Param("year")
		month := c.Param("month")
		daynum := c.Param("day")
		day := year + "/" + month + "/" + daynum
		if day == Today() {
			// Open()
			var scienceNews ScienceNewsGroup
			scienceNews.Get(day)
			var inTheaters InTheatersGroup
			inTheaters.Get(day)
			var groceryList CowyoGroup
			groceryList.Get(day, "GroceryList")
			var coolList CowyoGroup
			coolList.Get(day, "CoolList")
			var todoList CowyoGroup
			todoList.Get(day, "TodoList")
			var weather WeatherData
			weather.Get(day)
			// Close()

			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"Date":        day,
				"InTheaters":  inTheaters.Data,
				"GroceryList": groceryList.Data,
				"CoolList":    coolList.Data,
				"TodoList":    todoList.Data,
				"Weather":     weather,
			})
		} else {
			var p SelfData
			p.Get(day)
			c.Data(http.StatusOK, "text/html", []byte(p.Html))
		}

	})

	router.Run(":8001")

	// Exit 5 seconds later to let time for the request to be done.
	<-time.After(3 * time.Second)
	fmt.Println("HI")
}
