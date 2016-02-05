package main

import (
	"time"

	"github.com/celrenheit/spider"
	"github.com/celrenheit/spider/schedule"
)

func main() {
	// Create a new scheduler
	scheduler := spider.NewScheduler()

	// Register the spider to be scheduled every 15 seconds
	// Alternatively, you can choose a cron schedule
	// This will run every minute of every day
	// scheduler.Add(schedule.Cron("* * * * *"), LionelMessiSpider)
	//scheduler.Add(schedule.Cron("* * * * *"), MovieScoreSpider)
	scheduler.Add(schedule.Every(8*time.Second), ApNewsSpider)
	scheduler.Add(schedule.Every(5*time.Second), HackerNewsSpider)
	scheduler.Add(schedule.Every(7*time.Second), BestMoviesInTheatersSpider)

	// Start the scheduler
	scheduler.Start()

	// Exit 5 seconds later to let time for the request to be done.
	// Depends on your internet connection
	<-time.After(9 * time.Second)
}
