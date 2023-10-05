# ioc

Inversion of Control (IoC)

## Feature

* 1) Support service as singleton and transient

* 2) Support resolve service by parent if not found in current

* 3) Support inject to function or *struct with services that has registered

  Should add struct tag 'ioc-inject:"true"' to field if want to be injected, but field type `ioc.Resolver` is not necessary.

* 4) Support override exists service

  Register to parent's container, and then register to current's to override parent's.

* 5) Support inject to singleton instance automatically

  Inject to singleton instance and it's method `Initialized(XXX)` automatically.

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

type Interface3 interface {
    GetName() string
}

type Class3 struct {
    Name     string
}

func (c *Class3) GetName() string {
    return "Class3-" + c.Name
}

type Class4 struct {
    Name     string
}

func (c *Class4) GetName() string {
    return "Class3-" + c.Name
}

func main() {
    // register service to *struct
    ioc.AddSingleton[*Class2](&Class2{Name: "Jerry Bai"})
    ioc.AddTransient[*Class1](func() *Class1 {
        var svc Class1
        // inject to *struct
        ioc.Inject(&svc)
    }

    // register service to interface.
    ioc.AddSingleton[Interface2](&Class2{Name: "Jerry Bai"})
    ioc.AddTransient[Interface1](func() Interface1 {
        var svc Class1
        // inject to *struct
        ioc.Inject(&svc)
    }

    // get service from ioc
    c1 := ioc.GetService[*Class1]
    c2 := ioc.GetService[*Class2]
    i1 := ioc.GetService[Interface1]
    i2 := ioc.GetService[Interface2]

    // inject to function
    ioc.Inject(func(c1 *Class1, c2 *Class2, i1 Interface1, i2 Interface2, resolver ioc.Resolver) {
        println("c1.C2Name=", c1.C2.Name)
        println("c2.Name=", c2.Name)
        println("i1.GetC2Name=()", i1.GetC2Name())
        println("i2.GetName=()", i2.GetName())
    })

    // override exists service
    c := ioc.New()
    ioc.SetParent(c)
    ioc.AddSingletonToC[Interface3](c, &Class3{Name: "Jerry Bai"}) // add service to parent's container
    i3 := ioc.GetService[Interface3]() // *Class3, 'Interface3' only exists in parent's container
    ioc.AddSingleton[Interface3](&Class4{Name: "Jerry Bai"}) // add service to global's container
    i3 = ioc.GetService[Interface3]() // *Class4, 'Interface3' exists in both global and parent's container
}
```

## Benchmark

```sh
go test -run=none -count=1 -benchtime=1000000x -benchmem -bench=. ./...

goos: linux
goarch: amd64
pkg: github.com/berkaroad/ioc
cpu: AMD Ryzen 7 5800H with Radeon Graphics         
BenchmarkInjectToFunc-4          1000000               639.8 ns/op           128 B/op          5 allocs/op
BenchmarkInjectToStruct-4        1000000               463.4 ns/op            48 B/op          3 allocs/op
PASS
ok      github.com/berkaroad/ioc        1.107s
```

## Release Notes

### v1.1.1 (2023-10-05)

* 1) fix bug: replace exists service sometimes doesn't work

### v1.1 (2023-10-04)

* 1) add function `New() Container`

* 2) `Inject(target any)` can accept `reflect.Value`

* 3) `Resolve(serviceType reflect.Type) reflect.Value` will auto inject to singleton instance and it's method `Initialized(XXX)`

* 4) `SetParent(parent Resolver)` will append parent to exists parent

  can replace exists service, by register to parent's container and register to current's to override parent's.

* 5) add convenient functions for interface `Container`

  new functions `AddSingletonToC(XXX)`, `AddTransientToC(XXX)`, `GetServiceFromC(XXX)`, `InjectFromC(XXX)`.

  both last parent and new parent can also resolve services.

* 6) move `SetParent(parent Resolver)` from interface `Container` to `Resolver`

* 7) split method `Register(XXX)` to `RegisterSingleton(XXX)` and `RegisterTransient(XXX)`

* 8) 39% faster than v0.1.1

### v1.0 (2023-10-01)

refactor ioc: for simple and performance.

* 1) add convenient functions

* 2) support inject to function and *struct

  Should add struct tag 'ioc-inject:"true"' to field if want to be injected, but field type `ioc.Resolver` is not necessary.

* 3) 50% faster than v0.1.1, when injecting to function

  Compare with `Container.Invoke(f interface{}) ([]reflect.Value, error)` in v0.1.1

* 4) remove interface `Initializer`

  Because it is not necessary.

* 5) rename interface `ReadonlyContainer` to `Resolver`

  Just for `SetParent(parent Resolver)` to resolve by parent.

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
