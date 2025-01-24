package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/carsondecker/pokedex-cli/internal"
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

var cache *internal.Cache

func CreateCache(seconds int) {
	cache = internal.NewCache(time.Duration(seconds) * time.Second)
}

func getMapData(url string, mapConfig *Config) (MapData, error) {
	cachedData, ok := cache.Get(url)
	if ok {
		var decodedCache MapData
		if err := json.Unmarshal(cachedData, &decodedCache); err != nil {
			return MapData{}, errors.New("could not decode cache")
		}

		mapConfig.Next = decodedCache.Next
		mapConfig.Prev = decodedCache.Previous

		return decodedCache, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return MapData{}, errors.New("could not fetch data")
	}
	defer res.Body.Close()

	var resData MapData

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return MapData{}, errors.New("could not decode response for cache")
	}
	if err := json.Unmarshal(body, &resData); err != nil {
		return MapData{}, errors.New("could not decode response")
	}

	cache.Add(url, body)

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
