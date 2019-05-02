package spider

type ISpider interface {
	Do() error
	LoadCount() int32
	Name() string
}
