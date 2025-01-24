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

type config struct {
	next string
	prev string
}

var cache *internal.Cache

var mapConfig config = config{
	next: "https://pokeapi.co/api/v2/location-area/?offset=0&limit=0",
	prev: "",
}

func CreateCache(seconds int) {
	cache = internal.NewCache(time.Duration(seconds) * time.Second)
}

func getMapData(url string) (MapData, error) {
	cachedData, ok := cache.Get(url)
	if ok {
		var decodedCache MapData
		if err := json.Unmarshal(cachedData, &decodedCache); err != nil {
			return MapData{}, errors.New("could not decode cache")
		}

		mapConfig.next = decodedCache.Next
		mapConfig.prev = decodedCache.Previous

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

	mapConfig.next = resData.Next
	mapConfig.prev = resData.Previous

	return resData, nil
}

func GetNextMapData() (MapData, error) {
	if mapConfig.next == "" {
		return MapData{}, errors.New("no next pages")
	}
	return getMapData(mapConfig.next)
}

func GetPrevMapData() (MapData, error) {
	if mapConfig.prev == "" {
		return MapData{}, errors.New("no previous pages")
	}
	return getMapData(mapConfig.prev)
}
