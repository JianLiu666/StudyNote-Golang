package main

import (
	"errors"
	"log"
	"sync"

	"golang.org/x/sync/singleflight"
)

var g singleflight.Group

var errorNotExist = errors.New("not exist")

func main() {
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			data, err := getDataBySingleFlight("key")
			if err != nil {
				log.Print(err)
				return
			}
			log.Print(data)
		}()
	}

	wg.Wait()
}

func getDataByDefault(key string) (string, error) {
	data, err := getDataFromCache(key)

	if err == errorNotExist {
		data, err = getDataFromDatabase(key)
		if err != nil {
			log.Println(err)
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	return data, nil
}

func getDataBySingleFlight(key string) (string, error) {
	data, err := getDataFromCache(key)

	if err == errorNotExist {
		v, err, _ := g.Do(key, func() (any, error) {
			return getDataFromDatabase(key)
		})

		if err != nil {
			log.Println(err)
			return "", err
		}
		data = v.(string)

	} else if err != nil {
		return "", err
	}

	return data, nil
}

func getDataFromCache(key string) (string, error) {
	// 模擬 cache missed
	return "", errorNotExist
}

func getDataFromDatabase(key string) (string, error) {
	// 模擬從 DB 中查詢數據
	log.Printf("get %v from database", key)
	return "data", nil
}
