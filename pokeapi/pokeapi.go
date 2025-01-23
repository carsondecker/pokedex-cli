package pokeapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

type MapData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	Next string
	Prev string
}

func getMapData(url string, mapConfig *Config) (MapData, error) {
	res, err := http.Get(url)
	if err != nil {
		return MapData{}, errors.New("could not fetch data")
	}
	defer res.Body.Close()

	var resData MapData

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&resData); err != nil {
		return MapData{}, errors.New("could not decode response")
	}

	mapConfig.Next = resData.Next
	mapConfig.Prev = resData.Previous
	return resData, nil
}

func GetNextMapData(mapConfig *Config) (MapData, error) {
	if mapConfig.Next == "" {
		return MapData{}, errors.New("no next pages")
	}
	return getMapData(mapConfig.Next, mapConfig)
}

func GetPrevMapData(mapConfig *Config) (MapData, error) {
	if mapConfig.Prev == "" {
		return MapData{}, errors.New("no previous pages")
	}
	return getMapData(mapConfig.Prev, mapConfig)
}
