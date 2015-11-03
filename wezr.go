package main

import "gopkg.in/yaml.v2"
import "github.com/alexflint/go-arg"

import "encoding/json"
import "fmt"
import "io/ioutil"
import "log"
import "net/http"
import "os"
import "path/filepath"
import "strconv"

const VERSION = "0.1.0"
const BASE_URL = "https://api.forecast.io/forecast/"
const OPTIONS = "?exclude=minutely,hourly,daily"

type Weather struct {
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
	Timezone  string     `json:"timezone"`
	Offset    int        `json:"offset"`
	Currently *DataPoint `json:"currently"`
	Minutely  *DataBlock `json:"minutely"`
	Hourly    *DataBlock `json:"Hourly"`
	Daily     *DataBlock `json:"Daily"`
	Alerts    []Alert    `json:"Alerts"`
	Flags     Flags      `json:"Flags"`
}

type DataPoint struct {
	Time                int     `json:"time"`
	Summary             string  `json:"summary"`
	Icon                string  `json:"icon"`
	PrecipIntensity     float64 `json:"precipIntensity"`
	PrecipProbability   float64 `json:"precipProbability"`
	PrecipType          string  `json:"precipType"`
	Temperature         float64 `json:"temperature"`
	ApparentTemperature float64 `json:"apparentTemperature"`
	DewPoint            float64 `json:"dewPoint"`
	Humidity            float64 `json:"humidity"`
	WindSpeed           float64 `json:"windSpeed"`
	WindBearing         float64 `json:"windBearing"`
	Visibility          float64 `json:"visibility"`
	CloudCover          float64 `json:"cloudCover"`
	Pressure            float64 `json:"pressure"`
	Ozone               float64 `json:"ozone"`
}

func (dp DataPoint) String() string {
	return dp.Summary + " " + strconv.FormatFloat(dp.Temperature, 'f', 1, 32) + " (feels like " + strconv.FormatFloat(dp.ApparentTemperature, 'f', 1, 32) + ") precipitation chance " + strconv.FormatFloat(dp.PrecipProbability, 'f', 1, 32)
}

type DataBlock struct {
	Summary string      `json:"summary"`
	Icon    string      `json:"icon"`
	Data    []DataPoint `json:"data"`
}

type Alert struct {
	Title       string `json:"title"`
	Expires     int    `json:"expires"`
	Description string `json:"description"`
	Uri         string `json:"uri"`
}

type Flags struct {
	DarkskyUnavailable string   `json:"darksky-unavailable"`
	DarkskyStations    []string `json:"darksky-stations"`
	DatapointStations  []string `json:"datapoint-stations"`
	IsdStations        []string `json:"isd-stations"`
	LampStations       []string `json:"lamp-stations"`
	MetarStations      []string `json:"metar-stations"`
	MetnoStations      []string `json:"metno-stations"`
	Sources            []string `json:"sources"`
	Units              string   `json:"units"`
}

type Config struct {
	ApiKey string `yaml:"api_key"`
	Lat    string `yaml:"lat"`
	Long   string `yaml:"long"`
}

type Args struct {
	Version bool
}

func get_weather(api_key, lat, long string, not_metric bool) *Weather {
	coords := lat + "," + long
	var units string
	if not_metric {
		units = "&units=us"
	} else {
		units = "&units=si"
	}
	url := BASE_URL + api_key + "/" + coords + OPTIONS + units
	log.Printf(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var weather *Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		switch v := err.(type) {
		case *json.SyntaxError:
			fmt.Println(string(body[v.Offset-40 : v.Offset]))
		}
	}
	return weather
}

func main() {
	var args struct {
		CfgFile   string `arg:"--config,help:location of the configuration file - default: ~/.wezr.yml"`
		NotMetric bool   `arg:"--not-metric,help:don't use metric units"`
		Version   bool   `arg:"-v,help:show the current version"`
	}
	filename, _ := filepath.Abs(os.Getenv("HOME") + "/.wezr.yml")
	args.CfgFile = filename
	arg.MustParse(&args)
	if args.Version {
		fmt.Printf("wezr version %v", VERSION)
		os.Exit(0)
	}
	yamlFile, err := ioutil.ReadFile(args.CfgFile)
	if err != nil {
		log.Fatal("Error reading configuration file")
	}
	config := Config{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal("Error parsing configuration file")
	}
	weather := get_weather(config.ApiKey, config.Lat, config.Long, args.NotMetric)
	fmt.Println(weather.Currently)
}
