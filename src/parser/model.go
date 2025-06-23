package parser

type MapDistinct map[string]struct{}

func (m MapDistinct) Add(key string) {
	m[key] = struct{}{}
}
