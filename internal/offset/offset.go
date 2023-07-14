package offset

import "sync"

type Offsets struct {
	L    sync.RWMutex
	data map[string]int
}

func New() *Offsets {
	return &Offsets{
		data: make(map[string]int),
	}
}

func (o *Offsets) Get(key string) int {
	return o.data[key]
}

func (o *Offsets) Set(key string, offset int) {
	o.data[key] = offset
}
