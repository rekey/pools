package pools

import (
	"sync"
)

type Fn func() error

type Pools struct {
	wg          sync.WaitGroup
	rw          sync.RWMutex
	fns         []Fn
	max         int
	thread      int
	canRun      bool
	stopOnError bool
}

func NewPools(max int, stopOnError bool) *Pools {
	if max < 0 {
		max = 1
	}
	return &Pools{
		max:         max,
		stopOnError: stopOnError,
	}
}

func (that *Pools) Push(fn Fn) {
	that.fns = append(that.fns, fn)
}

func (that *Pools) add() {
	that.rw.Lock()
	defer that.rw.Unlock()
	that.thread += 1
}

func (that *Pools) remove() {
	that.rw.Lock()
	defer that.rw.Unlock()
	that.thread -= 1
}

func (that *Pools) unshift() *Fn {
	that.rw.Lock()
	defer that.rw.Unlock()
	length := len(that.fns)
	if length == 0 {
		return nil
	}
	fn := that.fns[0]
	if length == 1 {
		that.fns = []Fn{}
	} else {
		that.fns = that.fns[1:length]
	}
	return &fn
}

func (that *Pools) Run() error {
	that.canRun = true
	var err error
	for that.canRun {
		if that.stopOnError && err != nil {
			break
		}
		if len(that.fns) > 0 && that.thread <= that.max {
			fn := that.unshift()
			if fn != nil {
				that.add()
				that.wg.Add(1)
				go func() {
					defer func() {
						that.remove()
						that.wg.Done()
					}()
					fnErr := (*fn)()
					if fnErr != nil {
						err = fnErr
					}
				}()
			}
		}
		if len(that.fns) == 0 {
			that.canRun = false
		}
	}
	that.wg.Wait()
	return err
}
