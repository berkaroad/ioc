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

  Inject to singleton instance and it's initialize method `Initialize(XXX)` or another one which the returns of method `InitializeMethodName() string` automatically.

  It will use zero value instead of panic if depended service not registerd.

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
    resolver ioc.Resolver
}

func (c *Class2) GetName() string {
    return c.Name
}

func (c *Class2) Initialize(resolver ioc.Resolver) string {
    c.resolver = resolver
    return c.Name
}

type Interface3 interface {
    GetName() string
}

type Class3 struct {
    Name     string
    resolver ioc.Resolver
}

func (c *Class3) GetName() string {
    return "Class3-" + c.Name
}

// specific custom initialize method name
func (c *Class3) InitializeMethodName() string {
    return "MyInitialize"
}

// custom initialize method
func (c *Class2) MyInitialize(resolver ioc.Resolver) string {
    c.resolver = resolver
    return c.Name
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
BenchmarkGetSingletonService-4           1000000                26.16 ns/op            0 B/op          0 allocs/op
BenchmarkGetTransientService-4           1000000               370.9 ns/op            48 B/op          1 allocs/op
BenchmarkGetTransientServiceNative-4     1000000               131.9 ns/op            48 B/op          1 allocs/op
BenchmarkInjectToFunc-4                  1000000               659.5 ns/op           144 B/op          5 allocs/op
BenchmarkInjectToFuncNative-4            1000000                89.26 ns/op            0 B/op          0 allocs/op
BenchmarkInjectToStruct-4                1000000               311.7 ns/op             0 B/op          0 allocs/op
BenchmarkInjectToStructNative-4          1000000                87.64 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/berkaroad/ioc        1.686s
```
