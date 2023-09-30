// Package ioc is Inversion of Control (IoC).
// Support singleton and transient.
//
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
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var globalContainer Container = New()
var resolverType reflect.Type = reflect.TypeOf((*Resolver)(nil)).Elem()

// New ioc container, and add singleton service 'ioc.Resolver' to it.
func New() Container {
	var c Container = &defaultContainer{}
	c.Register(reflect.TypeOf((*Resolver)(nil)).Elem(), false, c, nil)
	return c
}

// Inversion of Control container.
type Container interface {
	Resolver

	// Set parent resolver, for resolving from parent if service not found in current.
	SetParent(parent Resolver)

	// Register to add singleton instance or transient instance factory.
	//
	//  // service
	//  type Service1 interface {
	//      Method1()
	//  }
	//  // implementation of service
	//  type ServiceImplementation1 struct {
	//      Field1 string
	//  }
	//  func(si *ServiceImplementation1) Method1() {}
	//
	//  var container ioc.Container
	//  // interface as service, register as singleton
	//  err := container.Register(reflect.TypeOf((*Service1)(nil)).Elem(), true, &ServiceImplementation1{Field1: "abc"}, nil)
	//  // or *struct as service, register as singleton
	//  err = container.Register(reflect.TypeOf((*ServiceImplementation1)(nil)), true, &ServiceImplementation1{Field1: "abc"}, nil)
	//
	//  // interface as service, register as transient
	//  err = container.Register(reflect.TypeOf((*Service1)(nil)).Elem(), true, nil, func() Service1 {
	//	    return &ServiceImplementation1{Field1: "abc"}
	//  })
	//  // or *struct as service, register as transient
	//  err = container.Register(reflect.TypeOf((*ServiceImplementation1)(nil)), true, nil, func() *ServiceImplementation1 {
	//	    return &ServiceImplementation1{Field1: "abc"}
	//  })
	Register(serviceType reflect.Type, canOverride bool, singletonInstance any, transientInstanceFactory any) error
}

// Resolver can resolve service.
type Resolver interface {
	// Resolve to get service.
	//
	//	// service
	//	type Service1 interface {
	//	    Method1()
	//	}
	//	// implementation of service
	//	type ServiceImplementation1 struct {
	//	    Field1 string
	//	}
	//	func(si *ServiceImplementation1) Method1() {}
	//
	//	var container ioc.Container
	//	// interface as service
	//	service1 := container.Resolve(reflect.TypeOf((*Service1)(nil)).Elem())
	//	// or *struct as service
	//	service2 := container.Resolve(reflect.TypeOf((*ServiceImplementation1)(nil)))
	Resolve(serviceType reflect.Type) reflect.Value
}

// AddSingleton to add singleton instance.
// can override exists service if has set 'canOverride' to 'true' before.
//
// It will panic if 'TService' or 'instance' is invalid.
//
//	// service
//	type Service1 interface {
//	    Method1()
//	}
//	// implementation of service
//	type ServiceImplementation1 struct {
//	    Field1 string
//	}
//	func(si *ServiceImplementation1) Method1() {}
//
//	// interface as service
//	ioc.AddSingleton[Service1](true, &ServiceImplementation1{Field1: "abc"})
//	// or *struct as service
//	ioc.AddSingleton[*ServiceImplementation1](true, &ServiceImplementation1{Field1: "abc"})
func AddSingleton[TService any](canOverride bool, instance TService) {
	err := globalContainer.Register(reflect.TypeOf((*TService)(nil)).Elem(), canOverride, instance, nil)
	if err != nil {
		panic(err)
	}
}

// AddTransient to add transient service instance factory.
// can override exists service if has set 'canOverride' to 'true' before.
//
// It will panic if 'TService' or 'instance' is invalid.
//
//	 // service
//	 type Service1 interface {
//	     Method1()
//	 }
//	 // implementation of service
//	 type ServiceImplementation1 struct {
//	     Field1 string
//	 }
//	 func(si *ServiceImplementation1) Method1() {}
//
//	 // interface as service
//	 ioc.AddTransient[Service1](true, func() Service1 {
//		    return &ServiceImplementation1{Field1: "abc"}
//	 })
//	 // or *struct as service
//	 ioc.AddTransient[*ServiceImplementation1](true, func() *ServiceImplementation1 {
//		    return &ServiceImplementation1{Field1: "abc"}
//	 })
func AddTransient[TService any](canOverride bool, instanceFactory func() TService) {
	err := globalContainer.Register(reflect.TypeOf((*TService)(nil)).Elem(), canOverride, nil, instanceFactory)
	if err != nil {
		panic(err)
	}
}

// GetService to get service.
//
//	// service
//	type Service1 interface {
//	    Method1()
//	}
//	// implementation of service
//	type ServiceImplementation1 struct {
//	    Field1 string
//	}
//	func(si *ServiceImplementation1) Method1() {}
//
//	// interface as service
//	service1 := ioc.GetService[Service1]()
//	// or *struct as service
//	service2 := ioc.GetService[*ServiceImplementation1]()
func GetService[TService any]() TService {
	var instance TService
	instanceVal := globalContainer.Resolve(reflect.TypeOf((*TService)(nil)).Elem())
	if !instanceVal.IsValid() {
		return instance
	}
	instanceInterface := instanceVal.Interface()
	if instanceInterface != nil {
		if val, ok := instanceInterface.(TService); ok {
			instance = val
		}
	}
	return instance
}

// Inject to func or *struct with service in ioc container.
// Field with type 'ioc.Resolver', will always been injected.
//
// It will panic if param type in func not registered in ioc container.
//
//	// service
//	type Service1 interface {
//	    Method1()
//	}
//
//	// implementation of service
//	type ServiceImplementation1 struct {
//	    Field1 string
//	}
//	func(si *ServiceImplementation1) Method1() {}
//
//	// client
//	type Client struct {
//	    Field1 Service1 `ioc-inject:"true"`
//	    Field2 *ServiceImplementation1 `ioc-inject:"true"`
//	}
//	func(c *Client) Method1(p1 Service1, p2 *ServiceImplementation1) {
//	    c.Field1 = p1
//	    c.Field2 = p2
//	}
//
//	var c client
//	// inject to func
//	ioc.Inject(c.Method1)
//	// inject to *struct
//	ioc.Inject(&c)
func Inject(target any) {
	targetVal := reflect.ValueOf(target)
	targetType := reflect.TypeOf(target)
	if targetType.Kind() == reflect.Func {
		// inject to func
		var in = make([]reflect.Value, targetType.NumIn())
		for i := 0; i < targetType.NumIn(); i++ {
			argType := targetType.In(i)
			val := globalContainer.Resolve(argType)
			if !val.IsValid() {
				panic(fmt.Errorf("service '%v' not found in ioc container, when injecting to func", argType))
			} else {
				in[i] = val
			}
		}
		targetVal.Call(in)
	} else if targetType.Kind() == reflect.Pointer && targetType.Elem().Kind() == reflect.Struct {
		// inject to *struct
		structType := targetType.Elem()
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			canInject := field.Type == resolverType
			if !canInject {
				if val, ok := field.Tag.Lookup("ioc-inject"); ok && val == "true" {
					canInject = true
				}
			}
			if canInject {
				fieldVal := globalContainer.Resolve(field.Type)
				if fieldVal.IsValid() {
					targetVal.Elem().Field(i).Set(fieldVal)
				}
			}
		}
	}
}

// Set parent resolver, for resolving from parent if service not found in current.
func SetParent(parent Resolver) {
	globalContainer.SetParent(parent)
}

var _ Container = (*defaultContainer)(nil)

type defaultContainer struct {
	bindings sync.Map
	parent   Resolver
}

func (c *defaultContainer) Resolve(serviceType reflect.Type) reflect.Value {
	if bindingVal, ok := c.bindings.Load(serviceType); ok {
		binding := bindingVal.(serviceBinding)
		if binding.Instance.IsValid() {
			return binding.Instance
		}
		return binding.InstanceFactory.Call(nil)[0]
	} else {
		parent := c.parent
		if parent != nil {
			return parent.Resolve(serviceType)
		} else {
			return reflect.Value{}
		}
	}
}

func (c *defaultContainer) SetParent(parent Resolver) {
	c.parent = parent
}

func (c *defaultContainer) Register(serviceType reflect.Type, canOverride bool, singletonInstance any, transinetInstanceFactory any) error {
	if singletonInstance == nil && transinetInstanceFactory == nil {
		return errors.New("both instance and instanceFactory are null")
	}
	if singletonInstance != nil {
		return c.addBinding(serviceBinding{ServiceType: serviceType, CanOverride: canOverride, Instance: reflect.ValueOf(singletonInstance)})
	} else {
		return c.addBinding(serviceBinding{ServiceType: serviceType, CanOverride: canOverride, InstanceFactory: reflect.ValueOf(transinetInstanceFactory)})
	}
}

func (c *defaultContainer) addBinding(binding serviceBinding) error {
	if binding.ServiceType == nil {
		return errors.New("type of service is null")
	}
	if binding.ServiceType.Kind() != reflect.Interface &&
		!(binding.ServiceType.Kind() == reflect.Pointer && binding.ServiceType.Elem().Kind() == reflect.Struct) {
		return fmt.Errorf("type of service '%v' should be an interface or *struct", binding.ServiceType)
	}
	if binding.Instance.IsValid() {
		if !binding.Instance.Type().AssignableTo(binding.ServiceType) {
			return fmt.Errorf("instance should implement the service '%v'", binding.ServiceType)
		}
	} else if binding.InstanceFactory.IsValid() {
		instanceFactoryType := binding.InstanceFactory.Type()
		if instanceFactoryType.Kind() != reflect.Func ||
			instanceFactoryType.NumIn() != 0 || instanceFactoryType.NumOut() != 1 ||
			!instanceFactoryType.Out(0).AssignableTo(binding.ServiceType) {
			return fmt.Errorf("type of instanceFactory should be a func with no params and return service '%v'", binding.ServiceType)
		}
	}
	if actual, loaded := c.bindings.LoadOrStore(binding.ServiceType, binding); loaded {
		existsBinding := actual.(serviceBinding)
		if existsBinding.CanOverride {
			c.bindings.Store(binding.ServiceType, binding)
		}
	}
	return nil
}

type serviceBinding struct {
	ServiceType     reflect.Type
	CanOverride     bool
	Instance        reflect.Value
	InstanceFactory reflect.Value
}
