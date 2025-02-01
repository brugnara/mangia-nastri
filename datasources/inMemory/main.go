package inMemory

import (
	"fmt"
	"mangia_nastri/datasources"
	"mangia_nastri/logger"
)

var log = logger.New("inMemory")

type InMemoryDataSource struct {
	data map[datasources.Hash]string
}

func New() *InMemoryDataSource {
	return &InMemoryDataSource{
		data: make(map[datasources.Hash]string),
	}
}

func (ds *InMemoryDataSource) Set(key datasources.Hash, value string) error {
	if key == "" {
		log.Printf("SET: key cannot be empty")
		return fmt.Errorf("key cannot be empty")
	}

	log.Printf("SET: Setting value %v for key %v", value, key)

	ds.data[key] = value
	return nil
}

func (ds *InMemoryDataSource) Get(key datasources.Hash) (string, error) {
	if key == "" {
		log.Printf("GET: key cannot be empty")
		return "", fmt.Errorf("key cannot be empty")
	}

	log.Printf("GET: Getting value for key %v", key)

	if _, ok := ds.data[key]; !ok {
		return "", fmt.Errorf("key not found")
	}

	return ds.data[key], nil
}
