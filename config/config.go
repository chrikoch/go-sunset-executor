package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//Config is a whole config
type Config struct {
	Entries []Entry `json:"entries"`
}

//Entry is one entry of a config
type Entry struct {
	Loc               Location `json:"location"`
	TimeOffsetMinutes int      `json:"offsetMinutes"`
	Target            string   `json:"target"`
	Method            string   `json:"method"`
}

//Location is the location, for which sunset/sunrise are computed
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

//ReadFromFile reads file filename into Config struct
func (c *Config) ReadFromFile(filename string) error {
	file, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)

	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &c)

	if err != nil {
		return err
	}

	return nil
}
