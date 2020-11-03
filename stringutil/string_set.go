package stringutil

type StringSet struct {
	values map[string]bool
}

func NewStringSet() *StringSet {
	return &StringSet{values: make(map[string]bool)}
}

func (ss *StringSet) Add(value ...string) {
	for _, val := range value {
		ss.values[val] = true
	}
}

func (ss *StringSet) Values() []string {
	result := make([]string, 0, len(ss.values))
	for k := range ss.values {
		result = append(result, k)
	}
	return result
}
