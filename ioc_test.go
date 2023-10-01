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
	"reflect"
	"testing"
)

func TestAddSingleton(t *testing.T) {
	t.Run("use interface as service and get service success", func(t *testing.T) {
		globalContainer = New()
		svc1 := &serviceInstance1{name: "instance1"}
		AddSingleton[service1](false, svc1)
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
		AddSingleton[*serviceInstance1](false, svc1)
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
				}
			}()
			AddSingleton[serviceInstance1](false, serviceInstance1{})
		}()
	})

	t.Run("override service that mark canOverride as true should success", func(t *testing.T) {
		globalContainer = New()
		svc1 := &serviceInstance1{name: "instance1"}
		AddSingleton[service1](true, svc1)
		svc1FromIoc := GetService[service1]()
		if svc1FromIoc == nil {
			t.Error("get service null")
			return
		}
		if svc1FromIoc != svc1 {
			t.Error("service should be singleton")
			return
		}

		AddTransient[service1](false, func() service1 { return &serviceInstance1{name: "override-instance1"} })
		svc1FromIoc = GetService[service1]()
		if svc1FromIoc == nil {
			t.Error("get service null")
			return
		}
		if svc1FromIoc == svc1 {
			t.Error("service should be another one")
			return
		}
		svc1FromIocAgain := GetService[service1]()
		if svc1FromIocAgain == nil {
			t.Error("get service null")
			return
		}
		if svc1FromIoc == svc1FromIocAgain {
			t.Error("service should be transient")
			return
		}
	})

	t.Run("override service that mark canOverride as false should fail", func(t *testing.T) {
		globalContainer = New()
		svc1 := &serviceInstance1{name: "instance1"}
		AddSingleton[service1](false, svc1)
		svc1FromIoc := GetService[service1]()
		if svc1FromIoc == nil {
			t.Error("get service null")
			return
		}
		if svc1FromIoc != svc1 {
			t.Error("service should be singleton")
			return
		}

		AddTransient[service1](true, func() service1 { return &serviceInstance1{name: "override-instance1"} })
		svc1FromIoc = GetService[service1]()
		if svc1FromIoc == nil {
			t.Error("get service null")
			return
		}
		if svc1FromIoc != svc1 {
			t.Error("service should be singleton")
			return
		}
	})
}

func TestAddTransient(t *testing.T) {
	t.Run("use interface as service and get service success", func(t *testing.T) {
		globalContainer = New()
		AddTransient[service2](false, func() service2 { return &serviceInstance2{name: "instance2"} })
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
		AddTransient[*serviceInstance2](false, func() *serviceInstance2 { return &serviceInstance2{name: "instance2"} })
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
				}
			}()
			AddTransient[serviceInstance1](false, func() serviceInstance1 { return serviceInstance1{} })
		}()
	})

	t.Run("override service that mark canOverride as true should success", func(t *testing.T) {
		globalContainer = New()
		AddTransient[service2](true, func() service2 { return &serviceInstance2{name: "instance2"} })
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

		AddSingleton[service2](false, &serviceInstance2{name: "override-instance2"})
		svc2FromIoc = GetService[service2]()
		if svc2FromIoc == nil {
			t.Error("get service null")
			return
		}
		svc2FromIocAgain := GetService[service2]()
		if svc2FromIocAgain == nil {
			t.Error("get service null")
			return
		}
		if svc2FromIoc != svc2FromIocAgain {
			t.Error("service should be singleton")
			return
		}
	})

	t.Run("override service that mark canOverride as false should fail", func(t *testing.T) {
		globalContainer = New()
		AddTransient[service2](false, func() service2 { return &serviceInstance2{name: "instance2"} })
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

		AddSingleton[service2](true, &serviceInstance2{name: "override-instance2"})
		svc2FromIoc = GetService[service2]()
		if svc2FromIoc == nil {
			t.Error("get service null")
			return
		}
		svc2FromIocAgain := GetService[service2]()
		if svc2FromIocAgain == nil {
			t.Error("get service null")
			return
		}
		if svc2FromIoc == svc2FromIocAgain {
			t.Error("service should be another one")
			return
		}
	})
}

func TestInject(t *testing.T) {
	t.Run("inject to func should success", func(t *testing.T) {
		globalContainer = New()
		AddSingleton[service3](false, &serviceInstance3{name: "instance3"})
		AddTransient[*serviceInstance3](false, func() *serviceInstance3 { return &serviceInstance3{name: "instance3"} })
		AddTransient[service4](false, func() service4 { return &serviceInstance4{name: "instance4"} })
		AddSingleton[*serviceInstance4](false, &serviceInstance4{name: "instance4"})

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
		AddSingleton[service3](false, &serviceInstance3{name: "instance3"})
		AddTransient[*serviceInstance3](false, func() *serviceInstance3 { return &serviceInstance3{name: "instance3"} })
		AddTransient[service4](false, func() service4 { return &serviceInstance4{name: "instance4"} })
		AddSingleton[*serviceInstance4](false, &serviceInstance4{name: "instance4"})

		var c client
		Inject(&c)

		if c.F3 != nil || c.F5 == nil || c.F5 == GetService[service4]() {
			t.Error("transient instance should same after inject, and only inject to field with tag 'ioc-inject:\"true\"'")
		}
		if c.F4 != nil || c.F6 == nil || c.F6 != GetService[*serviceInstance4]() {
			t.Error("singleton instance should different after inject, and only inject to field with tag 'ioc-inject:\"true\"'")
		}
	})
}

func TestSetParent(t *testing.T) {
	t.Run("can resolve from parent success", func(t *testing.T) {
		globalContainer = New()
		AddSingleton[service5](false, &serviceInstance5{name: "instance5"})
		AddTransient[*serviceInstance5](false, func() *serviceInstance5 { return &serviceInstance5{name: "instance5"} })

		anotherC := New()
		anotherC.Register(reflect.TypeOf((*service6)(nil)).Elem(), false, &serviceInstance6{name: "instance6"}, nil)
		anotherC.Register(reflect.TypeOf((*serviceInstance6)(nil)), false, nil, func() *serviceInstance6 { return &serviceInstance6{name: "instance6"} })

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
