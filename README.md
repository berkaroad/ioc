# ioc
Inversion of Control (IoC)

You can register a type as singleton or transient. Also you can register a type mapping to an interface as singleton or transient.


## Usage

    go get github.com/berkaroad/ioc


## Example

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
        C2Name string
    }

    func (c *Class1) InitFunc() interface{} {
        return func(c2 *Class2) {
            c.C2Name = c2.Name
        }
    }

    func (c *Class1) GetC2Name() string {
        return c.C2Name
    }

    type Class2 struct {
        Name string
    }

    func (c *Class2) InitFunc() interface{} {
        return func() {
            c.Name = "Tomcat"
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
        container.Invoke(func(c1 *Class1, c2 *Class2) {
            println("c1.C2Name=", c1.C2Name)
            println("c2.Name=", c2.Name)
        })

        // Like class's construction, inject interface instance
        container.Invoke(func(i1 Interface1, i2 Interface2) {
            println("i1.GetC2Name=()", i1.GetC2Name())
            println("i2.GetName=()", i2.GetName())
        })
    }