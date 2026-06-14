package datasources

type CacheableDataSource interface {
	Read(key string) ([]byte, error)
	Write(key string, data []byte) error
}
