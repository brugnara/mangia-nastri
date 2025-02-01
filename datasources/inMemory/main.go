package inMemory

import "mangia_nastri/datasources"

type InMemoryDataSource struct {
	data map[datasources.Hash]string
}

func New() *InMemoryDataSource {
	return &InMemoryDataSource{}
}

func (ds *InMemoryDataSource) Set(key datasources.Hash, value string) error {
	ds.data[key] = value
	return nil
}

func (ds *InMemoryDataSource) Get(key datasources.Hash) (string, error) {
	return ds.data[key], nil
}
