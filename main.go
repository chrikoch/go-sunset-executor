package main

import (
	"flag"
	"fmt"
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
	}
	fmt.Println(c)

	for _, e := range c.Entries {
		rise, set := sunrise.SunriseSunset(e.Loc.Latitude, e.Loc.Longitude, time.Now().Year(), time.Now().Month(), time.Now().Day())
		fmt.Println(rise)
		fmt.Println(set)
	}
}
