package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Celsius   float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Humidity  float64 `json:"humidity"`
	} `json:"main"`
	Sys struct {
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
}

type weatherData2 struct {
	CityName    string `json:"cityName"`
	TempDetails struct {
		Celsius   float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Humidity  float64 `json:"humidity"`
	} `json:"tempDetails"`
	Info struct {
		Country string `json:"country"`
		Sunrise string `json:"sunriseTime"`
		Sunset  string `json:"sunsetTime"`
	} `json:"nnfo"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}
	var apiData apiConfigData
	err = json.Unmarshal(bytes, &apiData)
	if err != nil {
		return apiConfigData{}, err
	}
	return apiData, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go! \n"))
}

func query(city string) (weatherData2, error) {
	appConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData2{}, err
	}

	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?&units=metric&appid=" + appConfig.OpenWeatherMapApiKey + "&q=" + city)

	if err != nil {
		return weatherData2{}, err
	}
	defer resp.Body.Close()
	var weather weatherData
	var finalWeather weatherData2
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return weatherData2{}, err
	}
	var t time.Time
	finalWeather.CityName = weather.Name
	finalWeather.TempDetails = weather.Main
	finalWeather.Info.Sunrise = t.Add(time.Duration(weather.Sys.Sunrise) * time.Millisecond).Format("15:04:05")
	finalWeather.Info.Sunset = t.Add(time.Duration(weather.Sys.Sunset) * time.Millisecond).Format("15:04:05")
	return finalWeather, nil
}

func main() {
	http.HandleFunc("/hello", hello)

	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})

	http.ListenAndServe(":8080", nil)
}
