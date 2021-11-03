# pools

## install

```shell
go get -v github.com/rekey/pools
```

## example

```go
//p := pools.NewPools(maxThreadCount, StopOnError)
p := pools.NewPools(10, true)
for i := 0; i < 15; i++ {
    (func(i int) {
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
```