package pool

import (
	"maps"
	"sync"
)

type Map[K comparable, V any] struct {
	p   sync.Pool
	cap int
	def map[K]V
}

func NewMap[K comparable, V any](capSave int, capCreate int, def map[K]V) *Map[K, V] {
	return &Map[K, V]{
		p: sync.Pool{
			New: func() any {
				c := capCreate
				if c == 0 {
					c = len(def)
				}
				m := make(map[K]V, c)
				maps.Copy(m, def)
				return m
			},
		},
		cap: capSave,
		def: def,
	}
}

func (p *Map[K, V]) Get() map[K]V {
	return p.p.Get().(map[K]V)
}

func (p *Map[K, V]) Save(m map[K]V) {
	if len(m) > p.cap {
		return
	}
	maps.Copy(m, p.def)
	maps.DeleteFunc(m, func(k K, v V) bool {
		_, ok := p.def[k]
		return !ok
	})
	p.p.Put(m)
}
