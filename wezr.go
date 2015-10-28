package main

import "gopkg.in/yaml.v2"

import "encoding/json"
import "fmt"
import "io/ioutil"
import "log"
import "net/http"
import "os"
import "path/filepath"
import "strconv"

const BASE_URL = "https://api.forecast.io/forecast/"
const OPTIONS = "?units=si&exclude=minutely,hourly,daily"

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

func get_weather(api_key, lat, long string) *Weather {
	coords := lat + "," + long
	url := BASE_URL + api_key + "/" + coords + OPTIONS
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

type Config struct {
	ApiKey string `yaml:"api_key"`
	Lat    string `yaml:"lat"`
	Long   string `yaml:"long"`
}

func main() {
	filename, _ := filepath.Abs(os.Getenv("HOME") + "/.wezr.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Error reading configuration file")
	}
	config := Config{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal("Error parsing configuration file")
	}
	weather := get_weather(config.ApiKey, config.Lat, config.Long)
	fmt.Println(weather.Currently)
}
