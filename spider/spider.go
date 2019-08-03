package spider

type Spider interface {
	Do() error
	LoadCount() int32
	Name() string
}
