package datasources

type DataSource interface {
	Set(key Hash, value Payload) error
	Get(key Hash) (Payload, error)
	Ready() <-chan bool
}
