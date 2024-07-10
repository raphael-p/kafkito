package utils

import (
	"sync"
)

type AutoIncrement struct {
	sync.Mutex
	id uint64
}

func (a *AutoIncrement) ID() uint64 {
	a.Lock()
	defer a.Unlock()

	a.id++
	return a.id
}
