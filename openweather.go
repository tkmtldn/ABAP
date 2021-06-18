package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const API_KEY = "60de2de0dd5236a5a4ad71357f028453"

type TemperatureInfo struct {
	NightTemp float64 `json:"night"`
	MornTemp  float64 `json:"morn"`
}

type WeatherDaily struct {
	Dt       int             `json:"dt"`
	Temp     TemperatureInfo `json:"temp"`
	Pressure int             `json:"pressure"`
}

type WeatherData struct {
	Daily [5]WeatherDaily `json:"daily"`
}

func getURL(lat, lon string) string {
	u, err := url.Parse("api.openweathermap.org")
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "https"
	u.Host = "api.openweathermap.org"
	u.Path = "/data/2.5/onecall"
	q := u.Query()
	q.Set("appid", API_KEY)
	q.Set("lat", lat)
	q.Set("lon", lon)
	q.Set("units", "metric")
	q.Set("exclude", "hourly,minutely")
	q.Set("cnt", "5")

	u.RawQuery = q.Encode()
	return u.String()
}

func timeParse(d int) string {
	day := strconv.Itoa(d)
	i, err := strconv.ParseInt(day, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm.Format("2006-01-02")
}

func differences(j WeatherData) (int, string, float64) {
	pressure := 0
	day := 0
	diff := 100.0

	for _, v := range j.Daily {
		if v.Pressure > pressure {
			pressure = v.Pressure
		}
		if math.Abs(v.Temp.NightTemp-v.Temp.MornTemp) < diff {
			diff = math.Abs(v.Temp.NightTemp - v.Temp.MornTemp)
			day = v.Dt
		}
	}

	return pressure, timeParse(day), diff
}

func main() {
	link := getURL("54.3333", "48.4")

	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}

	response, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var weather WeatherData
	err = json.Unmarshal(response, &weather)
	if err != nil {
		fmt.Println("error:", err)
	}

	pressure, date, temp := differences(weather)
	fmt.Printf("Maximum pressure forecast: %v hPa.\n", pressure)
	fmt.Printf("On %v we expect lowest difference between night and morning temperatures: %v degrees Celsius.\n", date, temp)
}
