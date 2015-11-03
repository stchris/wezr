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

```
Partly Cloudy 26.9 (feels like 29.2) precipitation chance 0.0
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
