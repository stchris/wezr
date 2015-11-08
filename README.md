# wezr
Shows weather info in the console. Weather data is provided by [forecast.io](https://forecast.io)

## Installation

* `go get github.com/stchris/wezr`


## Usage

* You need an API Key from [forecast.io](https://developer.forecast.io/)
* Find the coordinates you want weather data for
* `wezr` expects a `~/.wezr.yml` file with this structure:

```yaml
api_key: abcdef1234
lat: 12.345
long: 54.321
```
  
* Run `wezr`

## Sample output

```bash
$ wezr
Mostly Cloudy 8.3°C (feels like 6.8°C) precipitation chance 0%

$ wezr --units us
Mostly Cloudy 47.0°F (feels like 44.2°F) precipitation chance 0%

$ wezr --help
usage: wezr [--config CONFIG] [--not-metric] [--version] 

options:
  --config CONFIG        location of the configuration file - default: ~/.wezr.yml
  --units UNITS          display units: 'us' or 'si' (default)
  --version, -v          show the current version
```

## Creative usage

Add a cronjob (`crontab -e`) 

```bash
@hourly $GOHOME/bin/wezr > $HOME/.wezr.txt
```

and then use that info to greet you every time you open a new terminal, by putting this into your `.bash{rc,_profile}`:

```
echo "This is what it's like outside: `cat $HOME/.wezr.txt`"
```
