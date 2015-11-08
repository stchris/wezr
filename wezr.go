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
import "strings"

const VERSION = "0.2.0"
const BASE_URL = "https://api.forecast.io/forecast/"
const OPTIONS = "?exclude=minutely,hourly,daily"
const DEFAULT_TEMPLATE = "$summary $temp (feels like $apparentTemp) precipitation chance $precipitationChance"

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

// Renders the template string and replaces placeholders with values and units
func formatWeather(w *Weather, template, units string) string {
	crt := w.Currently
	// $temp
	var temp string
	if units == "si" {
		temp = fmt.Sprintf("%.1f°C", crt.Temperature)
	} else {
		temp = fmt.Sprintf("%.1f°F", crt.Temperature)
	}
	result := strings.Replace(template, "$temp", temp, -1)
	// $apparentTemp
	var apparentTemp string
	if units == "si" {
		apparentTemp = fmt.Sprintf("%.1f°C", crt.ApparentTemperature)
	} else {
		apparentTemp = fmt.Sprintf("%.1f°F", crt.ApparentTemperature)
	}
	result = strings.Replace(result, "$apparentTemp", apparentTemp, -1)
	// $precipitationChance
	precipitationChance := fmt.Sprintf("%d%%", int(crt.PrecipProbability*100))
	result = strings.Replace(result, "$precipitationChance", precipitationChance, -1)
	// $summary
	result = strings.Replace(result, "$summary", crt.Summary, -1)
	return result
}

func get_weather(api_key, lat, long string, units string) *Weather {
	coords := lat + "," + long
	url := BASE_URL + api_key + "/" + coords + OPTIONS + "&units=" + units
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
		CfgFile string `arg:"--config,help:location of the configuration file - default: ~/.wezr.yml"`
		Units   string `arg:"help:display units: 'us' or 'si' (default)"`
		Version bool   `arg:"-v,help:show the current version"`
	}
	filename, _ := filepath.Abs(os.Getenv("HOME") + "/.wezr.yml")
	args.CfgFile = filename
	args.Units = "si"
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
	units := "si"
	if args.Units != "" {
		units = args.Units
	}
	if units != "us" && units != "si" {
		log.Fatal("Unknown unit ", units)
	}
	weather := get_weather(config.ApiKey, config.Lat, config.Long, units)
	fmt.Println(formatWeather(weather, DEFAULT_TEMPLATE, units))
}
