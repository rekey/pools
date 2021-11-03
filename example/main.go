package main

import (
	"errors"
	"log"

	"github.com/rekey/pools"
)

func main() {
	p := pools.NewPools(10, true)
	for i := 0; i < 15; i++ {
		(func(i int) {
			//log.Println("p.Push", i)
			p.Push(func() error {
				if i == 3 {
					return errors.New("test error")
				}
				log.Println("p.Run", i)
				return nil
			})
		})(i)
	}
	err := p.Run()
	log.Println(err)
}
