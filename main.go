package main

import (
	"fmt"
	"time"

	"github.com/celrenheit/spider"
	"github.com/celrenheit/spider/schedule"
)

func main() {
	// Initialize all buckets

	Open()
	defer Close()

	// A Person struct consists of ID, Name, Age, Job.
	// var p ScienceNewsGroup
	// p.Data = []ScienceNews{}
	// p.Data = append(p.Data, ScienceNews{"Tidle1", "Summary1", "Url1"})
	// p.Data = append(p.Data, ScienceNews{"Title2", "Summary2", "Url2"})
	// p.save()
	// fmt.Println(p)

	var p ScienceNewsGroup
	err := p.Get(time.Now().String()[0:11])
	fmt.Println(p, err, len(p.Data))
	// check size of p.Data, if 0 then don't use

	// Create a new scheduler
	scheduler := spider.NewScheduler()

	// Register the spider to be scheduled every 15 seconds
	// Alternatively, you can choose a cron schedule
	// This will run every minute of every day
	// scheduler.Add(schedule.Cron("* * * * *"), LionelMessiSpider)
	//scheduler.Add(schedule.Cron("* * * * *"), MovieScoreSpider)
	// scheduler.Add(schedule.Every(8*time.Second), ApNewsSpider)
	// scheduler.Add(schedule.Every(5*time.Second), HackerNewsSpider)
	// scheduler.Add(schedule.Every(7*time.Second), BestMoviesInTheatersSpider)
	scheduler.Add(schedule.Every(2*time.Second), ScienceNewsSpider)

	// Start the scheduler
	scheduler.Start()

	// Exit 5 seconds later to let time for the request to be done.
	// Depends on your internet connection
	<-time.After(3 * time.Second)
}
