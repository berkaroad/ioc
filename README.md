# ioc

Inversion of Control (IoC)

## Feature

* 1) Support service as singleton and transient

* 2) Can resolve service by parent if not found in current

* 3) Can inject to function or *struct with services that has registered

  Should add struct tag 'ioc-inject:"true"' to field if want to be injected, but field type `ioc.Resolver` is not necessary.

## Usage

```go
package main

import (
    "github.com/berkaroad/ioc"
)

type Interface1 interface {
    GetC2Name() string
}

type Interface2 interface {
    GetName() string
}

type Class1 struct {
    Resolver ioc.Resolver
    C2       *Class2 `ioc-inject:"true"`
}

func (c *Class1) GetC2Name() string {
    return c.C2.Name
}

type Class2 struct {
    Name     string
}

func (c *Class2) GetName() string {
    return c.Name
}

func main() {
    // register service to *struct
    ioc.AddSingleton[*Class2](&Class2{Name: "Jerry Bai"})
    ioc.AddTransient[*Class1](func() *Class1 {
        var c1 Class1
        // inject to *struct
        ioc.Inject(&c1)
    }

    // register service to interface.
    ioc.AddSingleton[Interface2](&Class2{Name: "Jerry Bai"})
    ioc.AddTransient[Interface1](func() Interface1 {
        var c1 Class1
        // inject to *struct
        ioc.Inject(&c1)
    }

    // inject to function
    ioc.Inject(func(c1 *Class1, c2 *Class2, i1 Interface1, i2 Interface2, resolver ioc.Resolver) {
        println("c1.C2Name=", c1.C2.Name)
        println("c2.Name=", c2.Name)
        println("i1.GetC2Name=()", i1.GetC2Name())
        println("i2.GetName=()", i2.GetName())
    })
}
```

## Benchmark

```sh
go test -run=none -count=1 -benchtime=1000000x -benchmem -bench=. ./...

goos: linux
goarch: amd64
pkg: github.com/berkaroad/ioc
cpu: AMD Ryzen 7 5800H with Radeon Graphics         
BenchmarkInjectToFunc-4          1000000              1473 ns/op             128 B/op          5 allocs/op
BenchmarkInjectToStruct-4        1000000               862.4 ns/op            48 B/op          3 allocs/op
PASS
ok      github.com/berkaroad/ioc        2.348s
```

## Release Notes

### v1.0 (2023-10-01)

refactor ioc: for simple and performance.

* 1) add convenient functions

* 2) support inject to function and *struct

  Should add struct tag 'ioc-inject:"true"' to field if want to be injected, but field type `ioc.Resolver` is not necessary.

* 3) 50% faster than last version, when injecting to function

  Compare with `Container.Invoke()` in last version.

* 4) remove interface `Initializer`

  Because it is not necessary.

* 5) rename interface `ReadonlyContainer` to `Resolver`

  Just for `SetParent()` to resolve by parent.

* 6) simplify interface `Container`

  Can customize implementation just for compatibility with others

* 7) remove log

### v0.1.1 (2023-10-01)

* 1) add go mod base on 1.14

* 2) use benchmark

### v0.1 (2016-08-31)

* 1) Remove readonly lock,

* 2) Singleton instance's initialization called only once.

* 3) Performance is 15% faster than last version.
