package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nathan-osman/go-sunrise"

	"github.com/chrikoch/go-sunset-executor/config"
)

func main() {

	var configFilename string
	flag.StringVar(&configFilename, "config", "", "location of config-file")
	flag.Parse()

	if len(configFilename) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var c config.Config
	err := c.ReadFromFile(configFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, e := range c.Entries {
		go triggerExecutions(e)
	}

	//wait forever, because everything is done in the go-routines
	select {}
}

type date struct {
	day   int
	month time.Month
	year  int
}

func triggerExecutions(e config.Entry) {

	for {
		//wait until next execution time
		waitFor := time.Until(nextExecutionTime(e))
		log.Printf("have to wait %v seconds until sunset at %v, %v, offset %v\n", waitFor, e.Loc.Latitude, e.Loc.Longitude, e.TimeOffsetMinutes)
		time.Sleep(waitFor)

		//execute
		execute(e)
	}
}

func nextExecutionTime(e config.Entry) time.Time {
	var currentDate date
	currentDate.year, currentDate.month, currentDate.day = time.Now().Date()

	for {
		//currently only sunSET is supported
		_, sunsetCurrentDate := sunrise.SunriseSunset(e.Loc.Latitude, e.Loc.Longitude, currentDate.year, currentDate.month, currentDate.day)
		//add Offset
		executionTime := sunsetCurrentDate.Add(time.Duration(e.TimeOffsetMinutes) * time.Minute)

		//by adding 1 second to Now() we ensure the execution time isn't returned twice
		if executionTime.After(time.Now().Add(time.Second)) {
			return executionTime
		}

		//executionTime is in the past -> move to next day
		currentDate.year, currentDate.month, currentDate.day =
			time.Date(currentDate.year, currentDate.month, currentDate.day, 0, 0, 0, 0, time.Local).AddDate(0, 0, 1).Date()
	}
}

func execute(e config.Entry) {
	client := http.Client{}
	req, err := http.NewRequest(e.Method, e.Target, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("START calling %v with method %v\n", e.Target, e.Method)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	log.Printf("END calling %v with method %v\n", e.Target, e.Method)
	log.Printf("response Header: %v\n", resp.Header)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	log.Printf("response Body: %v\n", string(bodyBytes))

}
