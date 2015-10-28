# wezr
Shows weather info in the console. Weather data is provided by [forecast.io](https://forecast.io)

## Installation

* Clone this repository
* `go get`
* `go build`

## Usage

* You need an API Key from [forecast.io](https://developer.forecast.io/)
* Find the coordinates you want weather data for
* `wezr` expects a `~/.wezr.yml` file with this structure:

```yaml
api_key: abcdef1234
lat: 12.345
long: 54.321
```
  
* `./wezr`

