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

func getData[T any](url string) (T, error) {
	var emptyStruct T

	cachedData, ok := cache.Get(url)
	if ok {
		var decodedCache T
		if err := json.Unmarshal(cachedData, &decodedCache); err != nil {
			return emptyStruct, errors.New("could not decode cache")
		}

		return decodedCache, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return emptyStruct, errors.New("could not fetch data")
	}
	defer res.Body.Close()

	var resData T

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return emptyStruct, errors.New("could not decode response for cache")
	}
	if err := json.Unmarshal(body, &resData); err != nil {
		return emptyStruct, errors.New("could not decode response")
	}

	cache.Add(url, body)

	return resData, nil
}

func getMapData(url string) (MapData, error) {
	data, err := getData[MapData](url)
	if err != nil {
		return MapData{}, err
	}

	mapConfig.next = data.Next
	mapConfig.prev = data.Previous

	return data, nil
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
