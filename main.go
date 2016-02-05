package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	Open()
	defer Close()
	var p ScienceNewsGroup
	err := p.Get(time.Now().String()[0:10])
	fmt.Println(p, err, len(p.Data))

	// TODO check size of p.Data, if 0 then don't use

	// scheduler := spider.NewScheduler()
	// scheduler.Add(schedule.Every(2*time.Second), ScienceNewsSpider)
	// scheduler.Start()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"ScienceNews": p.Data,
		})
	})
	router.Run(":8080")

	// Exit 5 seconds later to let time for the request to be done.
	<-time.After(3 * time.Second)
	fmt.Println("HI")
}
