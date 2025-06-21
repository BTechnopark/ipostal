package parser

type MapDistinct map[string]struct{}

func (m MapDistinct) Add(key string) {
	m[key] = struct{}{}
}

type MapList map[string][]MapDistinct

func (m MapList) Add(key string, value MapDistinct) {
	m[key] = append(m[key], value)
}

var ListSort = []string{}
