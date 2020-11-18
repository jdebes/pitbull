package api

type RequestModel interface {
	Valid() error
}
