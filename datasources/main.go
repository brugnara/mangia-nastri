package datasources

type DataSource interface {
	Set(key Hash, value string) error
	Get(key Hash) (string, error)
}
