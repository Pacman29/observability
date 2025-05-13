package pool

import (
	"sync"
)

type Slice[V any] struct {
	p   sync.Pool
	cap int
	def []V
}

func NewSlice[V any](capSave int, capCreate int, def []V) *Slice[V] {
	return &Slice[V]{
		p: sync.Pool{
			New: func() any {
				c := capCreate
				if c == 0 {
					c = len(def)
				}
				s := make([]V, len(def), capCreate)
				copy(s, def)
				return s
			},
		},
		cap: capSave,
		def: def,
	}
}

func (p *Slice[V]) Get() []V {
	return p.p.Get().([]V)
}

func (p *Slice[V]) Save(s []V) {
	if len(s) > p.cap {
		return
	}
	s = s[:len(p.def)]
	copy(s, p.def)
	p.p.Put(s)
}
