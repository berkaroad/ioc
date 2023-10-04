// The MIT License (MIT)
//
// # Copyright (c) 2016 Jerry Bai
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package ioc

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddSingleton(t *testing.T) {
	t.Run("use interface as service and get service success", func(t *testing.T) {
		globalContainer = New()
		svc1 := &serviceInstance1{name: "instance1"}
		AddSingleton[service1](svc1)
		svc1FromIoc := GetService[service1]()
		if svc1FromIoc == nil {
			t.Error("get service null")
			return
		}
		if svc1FromIoc != svc1 {
			t.Error("service should be singleton")
			return
		}
	})

	t.Run("use *struct as service and get service success", func(t *testing.T) {
		globalContainer = New()
		svc1 := &serviceInstance1{name: "instance1"}
		AddSingleton[*serviceInstance1](svc1)
		svc1FromIoc := GetService[*serviceInstance1]()
		if svc1FromIoc == nil {
			t.Error("get service null")
			return
		}
		if svc1FromIoc != svc1 {
			t.Error("service should be singleton")
			return
		}
	})

	t.Run("invalid service should fail", func(t *testing.T) {
		globalContainer = New()
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Error("type of service 'serviceInstance1' should be interface or *struct")
				} else {
					fmt.Printf("panic: %v\n", r)
				}
			}()
			AddSingleton[serviceInstance1](serviceInstance1{})
		}()
	})

	t.Run("null service instance should fail", func(t *testing.T) {
		globalContainer = New()
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Error("instance couldn't be null")
				} else {
					fmt.Printf("panic: %v\n", r)
				}
			}()
			AddSingleton[*serviceInstance1](nil)
		}()
	})

	t.Run("use *struct as service and cycle reference in 'Initialize()' should fail", func(t *testing.T) {
		globalContainer = New()
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Error("cycle reference in 'Initialize()' in service 'serviceInstance9'")
				} else {
					fmt.Printf("panic: %v\n", r)
				}
			}()
			AddSingleton[*serviceInstance9](&serviceInstance9{})
		}()
	})

	t.Run("use interface as service and cycle reference in 'Initialize()' should fail", func(t *testing.T) {
		globalContainer = New()
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Error("cycle reference in 'Initialize()' in service 'serviceInstance9'")
				} else {
					fmt.Printf("panic: %v\n", r)
				}
			}()
			AddSingleton[service1](&serviceInstance10{})
		}()
	})
}

func TestAddTransient(t *testing.T) {
	t.Run("use interface as service and get service success", func(t *testing.T) {
		globalContainer = New()
		AddTransient[service2](func() service2 { return &serviceInstance2{name: "instance2"} })
		svc2FromIoc := GetService[service2]()
		if svc2FromIoc == nil {
			t.Error("get service null")
			return
		}
		if svc2FromIoc.GetName() != "instance2" {
			t.Error("name of service should be instance2")
			return
		}
		svc2 := svc2FromIoc
		svc2FromIoc = GetService[service2]()
		if svc2FromIoc == svc2 {
			t.Error("service should be transient")
			return
		}
	})

	t.Run("use *struct as service and get service success", func(t *testing.T) {
		globalContainer = New()
		AddTransient[*serviceInstance2](func() *serviceInstance2 { return &serviceInstance2{name: "instance2"} })
		svc2FromIoc := GetService[*serviceInstance2]()
		if svc2FromIoc == nil {
			t.Error("get service null")
			return
		}
		if svc2FromIoc.GetName() != "instance2" {
			t.Error("name of service should be instance2")
			return
		}
		svc2 := svc2FromIoc
		svc2FromIoc = GetService[*serviceInstance2]()
		if svc2FromIoc == svc2 {
			t.Error("service should be transient")
			return
		}
	})

	t.Run("invalid service should fail", func(t *testing.T) {
		globalContainer = New()
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Error("type of service 'serviceInstance1' should be interface or *struct")
				} else {
					fmt.Printf("panic: %v\n", r)
				}
			}()
			AddTransient[serviceInstance1](func() serviceInstance1 { return serviceInstance1{} })
		}()
	})

	t.Run("null service instance factory should fail", func(t *testing.T) {
		globalContainer = New()
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Error("instance factory couldn't be null")
				} else {
					fmt.Printf("panic: %v\n", r)
				}
			}()
			AddTransient[*serviceInstance1](nil)
		}()
	})
}

func TestGetService(t *testing.T) {
	t.Run("get service with func 'Initialize()' should success", func(t *testing.T) {
		globalContainer = New()
		AddSingleton[*serviceInstance7](&serviceInstance7{name: "instance7"})
		AddSingleton[*serviceInstance8](&serviceInstance8{})

		svc8 := GetService[*serviceInstance8]()
		if svc8.GetS7Name() != "instance7" {
			t.Error("should function 'initialize()' invoked success")
		}
	})

	t.Run("replace exists service should success", func(t *testing.T) {
		globalContainer = New()
		anotherC := New()
		SetParent(anotherC)

		AddSingletonToC[*serviceInstance7](anotherC, &serviceInstance7{name: "instance7"})
		AddSingletonToC[*serviceInstance8](anotherC, &serviceInstance8{})

		// replace exists service
		AddSingleton[*serviceInstance7](&serviceInstance7{name: "new-instance7"})

		svc8 := GetService[*serviceInstance8]()
		if svc8.GetS7Name() != "new-instance7" {
			t.Error("should replace exists service success")
		}
	})

	t.Run("func 'Initialize()' missing service should fail", func(t *testing.T) {
		globalContainer = New()
		AddSingleton[*serviceInstance8](&serviceInstance8{})
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Error("missing service should fail")
				} else {
					fmt.Printf("panic: %v\n", r)
				}
			}()
			svc8 := GetService[*serviceInstance8]()
			if svc8.GetS7Name() != "instance7" {
				t.Error("should function 'initialize()' invoked success")
			}
		}()
	})
}

func TestInject(t *testing.T) {
	t.Run("inject to func should success", func(t *testing.T) {
		globalContainer = New()
		AddSingleton[service3](&serviceInstance3{name: "instance3"})
		AddTransient[*serviceInstance3](func() *serviceInstance3 { return &serviceInstance3{name: "instance3"} })
		AddTransient[service4](func() service4 { return &serviceInstance4{name: "instance4"} })
		AddSingleton[*serviceInstance4](&serviceInstance4{name: "instance4"})

		var c client
		Inject(c.Func1)
		Inject((&c).Func1)
		if c.F1 == nil || c.F1 != GetService[service3]() {
			t.Error("singleton instance should same after inject")
		}
		if c.F2 == nil || c.F2 == GetService[*serviceInstance3]() {
			t.Error("transient instance should different after inject")
		}
		if c.F3 == nil || c.F3 == GetService[service4]() {
			t.Error("transient instance should different after inject")
		}
		if c.F4 == nil || c.F4 != GetService[*serviceInstance4]() {
			t.Error("singleton instance should same after inject")
		}
	})

	t.Run("inject to struct should success", func(t *testing.T) {
		globalContainer = New()
		AddSingleton[service3](&serviceInstance3{name: "instance3"})
		AddTransient[*serviceInstance3](func() *serviceInstance3 { return &serviceInstance3{name: "instance3"} })
		AddTransient[service4](func() service4 { return &serviceInstance4{name: "instance4"} })
		AddSingleton[*serviceInstance4](&serviceInstance4{name: "instance4"})

		var c client
		Inject(&c)

		if c.F3 != nil || c.F5 == nil || c.F5 == GetService[service4]() {
			t.Error("transient instance should same after inject, and only inject to field with tag 'ioc-inject:\"true\"'")
		}
		if c.F4 != nil || c.F6 == nil || c.F6 != GetService[*serviceInstance4]() {
			t.Error("singleton instance should different after inject, and only inject to field with tag 'ioc-inject:\"true\"'")
		}
	})

	t.Run("inject to impletementation of ioc.Resolve should ignore", func(t *testing.T) {
		globalContainer = New()
		c := &defaultContainer{}
		InjectFromC(c, c)
		if c.parent != nil {
			t.Error("inject to impletementation of ioc.Resolve should ignore")
		}
	})

	t.Run("inject to invalid reflect.Value should ignore", func(t *testing.T) {
		globalContainer = New()
		c := &defaultContainer{}
		InjectFromC(c, reflect.Value{})
	})

	t.Run("inject to null should ignore", func(t *testing.T) {
		globalContainer = New()
		c := &defaultContainer{}
		InjectFromC(c, (*serviceInstance1)(nil))
	})
}

func TestSetParent(t *testing.T) {
	t.Run("resolve from parent success", func(t *testing.T) {
		globalContainer = New()

		anotherC := New()
		AddSingletonToC[service6](anotherC, &serviceInstance6{name: "instance6"})
		AddTransientToC[*serviceInstance6](anotherC, func() *serviceInstance6 { return &serviceInstance6{name: "instance6"} })

		if svc := GetService[service6](); svc != nil {
			t.Error("service should not found in current")
		}
		if svc := GetService[*serviceInstance6](); svc != nil {
			t.Error("service should not found in current")
		}
		SetParent(anotherC)
		if svc := GetService[service6](); svc == nil {
			t.Error("service should found in parent")
		}
		if svc := GetService[*serviceInstance6](); svc == nil {
			t.Error("service should found in parent")
		}
	})

	t.Run("override parent's service success", func(t *testing.T) {
		globalContainer = New()

		anotherC := New()
		AddSingletonToC[service6](anotherC, &serviceInstance6{name: "instance6"})
		AddTransientToC[*serviceInstance6](anotherC, func() *serviceInstance6 { return &serviceInstance6{name: "instance6"} })
		SetParent(anotherC)
		if svc := GetService[service6](); svc == nil || svc.GetName() != "instance6" {
			t.Error("service should found in parent")
		}
		if svc := GetService[*serviceInstance6](); svc == nil || svc.GetName() != "instance6" {
			t.Error("service should found in parent")
		}

		AddSingleton[service6](&serviceInstance6{name: "override-instance6"})
		AddTransient[*serviceInstance6](func() *serviceInstance6 { return &serviceInstance6{name: "override-instance6"} })
		if svc := GetService[service6](); svc == nil || svc.GetName() != "override-instance6" {
			t.Error("service should override in global")
		}
		if svc := GetService[*serviceInstance6](); svc == nil || svc.GetName() != "override-instance6" {
			t.Error("service should override in global")
		}
	})

	t.Run("parent is null or equals with last one should ignore", func(t *testing.T) {
		globalContainer = New()

		anotherC := New()
		AddSingletonToC[service6](anotherC, &serviceInstance6{name: "instance6"})
		AddTransientToC[*serviceInstance6](anotherC, func() *serviceInstance6 { return &serviceInstance6{name: "instance6"} })

		if svc := GetService[service6](); svc != nil {
			t.Error("service should not found in current")
		}
		if svc := GetService[*serviceInstance6](); svc != nil {
			t.Error("service should not found in current")
		}
		SetParent(anotherC)
		if svc := GetService[service6](); svc == nil {
			t.Error("service should found in parent")
		}
		if svc := GetService[*serviceInstance6](); svc == nil {
			t.Error("service should found in parent")
		}
		SetParent(nil)
		if svc := GetService[service6](); svc == nil {
			t.Error("service should found in parent")
		}
		if svc := GetService[*serviceInstance6](); svc == nil {
			t.Error("service should found in parent")
		}
		SetParent(anotherC)
		if svc := GetService[service6](); svc == nil {
			t.Error("service should found in parent")
		}
		if svc := GetService[*serviceInstance6](); svc == nil {
			t.Error("service should found in parent")
		}
	})

	t.Run("set parent twice, last parent should been new parent's parent ", func(t *testing.T) {
		globalContainer = New()

		anotherC := New()
		AddSingletonToC[service6](anotherC, &serviceInstance6{name: "instance6"})
		AddTransientToC[*serviceInstance6](anotherC, func() *serviceInstance6 { return &serviceInstance6{name: "instance6"} })

		anotherC2 := New()
		AddSingletonToC[service5](anotherC2, &serviceInstance5{name: "instance5"})
		AddTransientToC[*serviceInstance5](anotherC2, func() *serviceInstance5 { return &serviceInstance5{name: "instance5"} })
		if svc := GetService[service6](); svc != nil {
			t.Error("service should not found in current")
		}
		if svc := GetService[*serviceInstance6](); svc != nil {
			t.Error("service should not found in current")
		}
		SetParent(anotherC)
		if svc := GetService[service6](); svc == nil {
			t.Error("service should found in parent")
		}
		if svc := GetService[*serviceInstance6](); svc == nil {
			t.Error("service should found in parent")
		}
		SetParent(anotherC2)
		if svc := GetService[service6](); svc == nil {
			t.Error("service should found in parent")
		}
		if svc := GetService[*serviceInstance6](); svc == nil {
			t.Error("service should found in parent")
		}
		if svc := GetService[service5](); svc == nil {
			t.Error("service should found in parent")
		}
		if svc := GetService[*serviceInstance5](); svc == nil {
			t.Error("service should found in parent")
		}
	})
}

func TestRegisterSingleton(t *testing.T) {
	t.Run("null service type should fail", func(t *testing.T) {
		globalContainer = New()

		c := New()
		err := c.RegisterSingleton(nil, nil)
		if err == nil {
			t.Error("service type should be null")
		}
	})

	t.Run("null service instance should fail", func(t *testing.T) {
		globalContainer = New()

		c := New()
		err := c.RegisterSingleton(reflect.TypeOf((*serviceInstance1)(nil)), nil)
		if err == nil {
			t.Error("null service instance should fail")
		}
	})

	t.Run("service instance should impletement service", func(t *testing.T) {
		globalContainer = New()

		c := New()
		err := c.RegisterSingleton(reflect.TypeOf((*service2)(nil)).Elem(), &serviceInstance1{})
		if err == nil {
			t.Error("service instance should impletement service")
		}
	})
}

func TestRegisterTransient(t *testing.T) {
	t.Run("", func(t *testing.T) {
		globalContainer = New()

		c := New()
		err := c.RegisterTransient(nil, nil)
		if err == nil {
			t.Error("null service type should fail")
		}
	})

	t.Run("null service instance factory should fail", func(t *testing.T) {
		globalContainer = New()

		c := New()
		err := c.RegisterTransient(reflect.TypeOf((*serviceInstance1)(nil)), nil)
		if err == nil {
			t.Error("null service instance factory should fail")
		}
	})

	t.Run("return value of service instance factory should impletement service", func(t *testing.T) {
		globalContainer = New()

		c := New()
		err := c.RegisterTransient(reflect.TypeOf((*service2)(nil)).Elem(), func() *serviceInstance1 {
			return &serviceInstance1{}
		})
		if err == nil {
			t.Error("service instance should impletement service")
		}
	})
}

type service1 interface {
	GetName() string
}

type service2 interface {
	GetName() string
	Rename(name string)
}

type service3 interface {
	service1
}

type service4 interface {
	service2
}

type service5 interface {
	service1
}

type service6 interface {
	service2
}

type serviceInstance1 struct {
	name string
}

func (instance *serviceInstance1) GetName() string {
	return instance.name
}

type serviceInstance2 struct {
	name string
}

func (instance *serviceInstance2) GetName() string {
	return instance.name
}

func (instance *serviceInstance2) Rename(name string) {
	instance.name = name
}

type serviceInstance3 struct {
	name string
}

func (instance *serviceInstance3) GetName() string {
	return instance.name
}

type serviceInstance4 struct {
	name string
}

func (instance *serviceInstance4) GetName() string {
	return instance.name
}

func (instance *serviceInstance4) Rename(name string) {
	instance.name = name
}

type serviceInstance5 struct {
	name string
}

func (instance *serviceInstance5) GetName() string {
	return instance.name
}

type serviceInstance6 struct {
	name string
}

func (instance *serviceInstance6) GetName() string {
	return instance.name
}

func (instance *serviceInstance6) Rename(name string) {
	instance.name = name
}

type client struct {
	F1 service3
	F2 *serviceInstance3
	F3 service4
	F4 *serviceInstance4
	F5 service4          `ioc-inject:"true" `
	F6 *serviceInstance4 `ioc-inject:"true" `
}

func (c *client) Func1(p1 service3, p2 *serviceInstance3, p3 service4, p4 *serviceInstance4) {
	c.F1 = p1
	c.F2 = p2
	c.F3 = p3
	c.F4 = p4
}

type serviceInstance7 struct {
	name string
}

func (instance *serviceInstance7) GetName() string {
	return instance.name
}

type serviceInstance8 struct {
	s7 *serviceInstance7
}

func (instance *serviceInstance8) GetS7Name() string {
	return instance.s7.name
}

func (instance *serviceInstance8) Initialize(s7 *serviceInstance7) {
	instance.s7 = s7
}

type serviceInstance9 struct {
	name string
	s9   *serviceInstance9
}

func (instance *serviceInstance9) GetName() string {
	return instance.name
}

func (instance *serviceInstance9) Initialize(s9 *serviceInstance9) {
	instance.s9 = s9
}

type serviceInstance10 struct {
	name string
	s10  service1
}

func (instance *serviceInstance10) GetName() string {
	return instance.name
}

func (instance *serviceInstance10) Initialize(s10 service1) {
	instance.s10 = s10
}
