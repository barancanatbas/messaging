package interfaces

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string, value interface{}) error
	Delete(key string) error
	LPush(key string, value interface{}) error
	LRange(key string, start, stop int64, data interface{}) error
	LIndex(key string, index int, data interface{}) error
}
