package stringutil

type StringIterator interface {
	Value() string
	Next() bool
	Err() error
}

type stringIterator struct {
	current    int
	data       []string
	err        error
	loadMoreFn func() ([]string, error)
}

func NewStringIterator(loadMoreFn func() ([]string, error)) StringIterator {
	return &stringIterator{
		current:    0,
		data:       []string{},
		loadMoreFn: loadMoreFn,
	}
}

func (i *stringIterator) Value() string {
	return i.data[i.current]
}

func (i *stringIterator) Next() bool {
	if i.current+1 >= len(i.data) {
		return i.loadMore()
	}

	i.current++
	return true
}

func (i *stringIterator) loadMore() bool {
	data, err := i.loadMoreFn()
	if err != nil {
		i.err = err
		return false
	}
	i.data = data
	i.current = 0

	return len(i.data) > 0
}

func (i *stringIterator) Err() error {
	return i.err
}
