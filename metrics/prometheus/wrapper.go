package prometheus

import "sync"

type pair[T any] struct {
	metric *T
	labels []string
}

type metricWrapper[T any] struct {
	lock *sync.RWMutex
	m    map[string]pair[T]
}

func newWrapper[T any](size int) *metricWrapper[T] {
	return &metricWrapper[T]{
		lock: &sync.RWMutex{},
		m:    make(map[string]pair[T], size),
	}
}

func (w *metricWrapper[T]) Add(key string, metric *T, labelNames []string) {
	w.lock.Lock()
	w.m[key] = pair[T]{
		metric: metric,
		labels: labelNames,
	}
	w.lock.Unlock()
}

func (w *metricWrapper[T]) Get(key string) (*T, []string, bool) {
	w.lock.RLock()
	p, ok := w.m[key]
	w.lock.RUnlock()
	if !ok {
		return nil, nil, false
	}
	return p.metric, p.labels, true
}
