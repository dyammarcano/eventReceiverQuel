package util

//import (
//	"fmt"
//	"github.com/robfig/cron/v3"
//	"time"
//)
//
//// http://godoc.org/github.com/robfig/cron
//
//func NextCronExecution(cronExpression string) (time.Time, error) {
//	c := cron.New()
//	//parsedCron, err := cron.ParseStandard(cronExpression)
//	//if err != nil {
//	//	return time.Time{}, err
//	//}
//	//
//	//now := time.Now()
//	//nextExecutionTime := parsedCron.Next(now)
//	//
//	//return nextExecutionTime, nil
//
//	_, err := c.AddFunc(cronExpression, func() {
//		fmt.Println("Executing cron job at 5:00 AM")
//	})
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	c.Start()
//}

/*
package main

import (
	"fmt"
	"github.com/robfig/cron"
	"time"
)

func main() {
	cronExpression := "0 5 * * *"
	nextExecutionTimeCh := make(chan time.Time)

	go NextCronExecution(cronExpression, nextExecutionTimeCh)

	nextExecutionTime := <-nextExecutionTimeCh
	fmt.Println("Next execution time:", nextExecutionTime)
}

func NextCronExecution(cronExpression string, nextExecutionTimeCh chan<- time.Time) {
	c := cron.New()
	parsedCron, err := cron.ParseStandard(cronExpression)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	now := time.Now()
	nextExecutionTime := parsedCron.Next(now)

	nextExecutionTimeCh <- nextExecutionTime
	close(nextExecutionTimeCh)

	_, err = c.AddFunc(cronExpression, func() {
		nextExecutionTime := parsedCron.Next(time.Now())
		nextExecutionTimeCh <- nextExecutionTime
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	c.Start()

	// Keep the program running
	select {}
}

*/
