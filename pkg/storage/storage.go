package storage

type Storage interface {
	Put(bucket, key string, value []byte) error
	Get(bucket, key string) ([]byte, error)
	Exist(bucket, key string) (bool, error)
}
