# ioc

Inversion of Control (IoC)

You can register a type as singleton or transient.
Also you can register a type mapping to an interface as singleton or transient.

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
    C2Name          string
    isInitialized   bool
}

func (c *Class1) InitFunc() interface{} {
    return func(c2 *Class2) {
        if !c.isInitialized {
            c.isInitialized = true
            c.C2Name = c2.Name
        }
    }
}

func (c *Class1) GetC2Name() string {
    return c.C2Name
}

type Class2 struct {
    Name            string
    isInitialized   bool
}

func (c *Class2) InitFunc() interface{} {
    return func() {
        if !c.isInitialized {
            c.isInitialized = true
            c.Name = "Tomcat"
        }
    }
}

func (c *Class2) GetName() string {
    return c.Name
}

func main() {
    var container = ioc.NewContainer()

    // Register class
    container.Register(&Class1{}, ioc.Singleton)
    container.Register(&Class2{Name: "Jerry Bai"}, ioc.Singleton)

    // Register class mapping to interface.
    container.RegisterTo(&Class1{}, (*Interface1)(nil), ioc.Transient)
    container.RegisterTo(&Class2{Name: "Jerry Bai"}, (*Interface2)(nil), ioc.Transient)

    // Like class's construction, inject class instance
    container.Invoke(func(c1 *Class1, c2 *Class2, roContainer ioc.ReadonlyContainer) {
        println("c1.C2Name=", c1.C2Name)
        println("c2.Name=", c2.Name)
    })

    // Like class's construction, inject interface instance
    container.Invoke(func(i1 Interface1, i2 Interface2, roContainer ioc.ReadonlyContainer) {
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
BenchmarkInjectToFunc-4          1000000              2251 ns/op             128 B/op          5 allocs/op
PASS
ok      github.com/berkaroad/ioc        2.260s
```

## Release Notes

### v0.1.1 (2023-10-01)

* 1) add go mod base on 1.14

* 2) use benchmark

### v0.1 (2016-08-31)

* 1) Remove readonly lock,

* 2) Singleton instance's initialization called only once.

* 3) Performance is 15% faster than last version.
