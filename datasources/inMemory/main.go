package inMemory

import (
	"fmt"
	"mangia_nastri/datasources"
	"mangia_nastri/logger"
)

type InMemoryDataSource struct {
	data map[datasources.Hash]string
	log  logger.Logger
}

func New(log *logger.Logger) *InMemoryDataSource {
	return &InMemoryDataSource{
		data: make(map[datasources.Hash]string),
		log:  log.CloneWithPrefix("inMemory"),
	}
}

func (ds *InMemoryDataSource) Set(key datasources.Hash, value string) error {
	if key == "" {
		ds.log.Printf("SET: key cannot be empty")
		return fmt.Errorf("key cannot be empty")
	}

	ds.log.Printf("SET: Setting value %v for key %v", value, key)

	ds.data[key] = value
	return nil
}

func (ds *InMemoryDataSource) Get(key datasources.Hash) (string, error) {
	if key == "" {
		ds.log.Printf("GET: key cannot be empty")
		return "", fmt.Errorf("key cannot be empty")
	}

	ds.log.Printf("GET: Getting value for key %v", key)

	if _, ok := ds.data[key]; !ok {
		return "", fmt.Errorf("key not found")
	}

	return ds.data[key], nil
}
