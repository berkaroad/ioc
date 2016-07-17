/*
Package ioc is Inversion of Control (IoC).
Support lifecycle, such as singleton and transient.
*/
package ioc

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
	"time"
)

var consoleLog = log.New(os.Stdout, "[ioc] ", log.LstdFlags)

// DEBUG is a switcher for debug
var DEBUG = false

// Lifecycle is singleton or transient
type Lifecycle int

const (
	// Singleton is single instance
	Singleton Lifecycle = 0
	// Transient is temporary instance
	Transient Lifecycle = 1
)

// Initializer is to init a struct.
type Initializer interface {
	InitFunc() interface{}
}

// ReadonlyContainer is a readonly container
type ReadonlyContainer interface {
	// Resolve is to get instance of type.
	Resolve(typ reflect.Type) reflect.Value
	// Invoke is to inject to function's params, such as construction.
	Invoke(f interface{}) ([]reflect.Value, error)
}

// Container is a container for ioc.
type Container interface {
	Initializer
	ReadonlyContainer
	// Register is to register a type as singleton or transient.
	Register(val interface{}, lifecycle Lifecycle)
	// RegisterTo is to register a interface as singleton or transient.
	RegisterTo(val interface{}, ifacePtr interface{}, lifecycle Lifecycle)
	// SetParent is to resolve parent's container if current hasn't registered a type.
	SetParent(parent ReadonlyContainer)
}

// NewContainer : create iocContainer
func NewContainer() Container {
	container := &iocContainer{}
	reflect.ValueOf(container.InitFunc()).Call(nil)
	container.RegisterTo(container, (*ReadonlyContainer)(nil), Singleton)
	return container
}

type iocContainer struct {
	isInitialized bool
	locker        *sync.RWMutex
	singleton     *singletonContainer
	transient     *transientContainer
	parent        ReadonlyContainer
}

func (container *iocContainer) InitFunc() interface{} {
	return func() {
		if !container.isInitialized {
			container.locker = &sync.RWMutex{}
			container.singleton = &singletonContainer{valuemapper: make(map[reflect.Type]reflect.Value)}
			container.transient = &transientContainer{typemapper: make(map[reflect.Type]reflect.Type)}
			container.isInitialized = true
		}
	}
}

func (container *iocContainer) Register(val interface{}, lifecycle Lifecycle) {
	defer container.locker.Unlock()
	container.locker.Lock()

	typ := reflect.TypeOf(val)
	if container.transient.Get(typ).IsValid() {
		lifecycle = Transient
	} else if container.singleton.Get(typ).IsValid() {
		lifecycle = Singleton
	}
	switch lifecycle {
	case Transient:
		container.transient.Map(val)
		if DEBUG {
			consoleLog.Printf("ioc.Container.Register transient '%v'.\n", typ)
		}
	case Singleton:
		container.singleton.Map(val)
		if DEBUG {
			consoleLog.Printf("ioc.Container.Register singleton '%v'.\n", typ)
		}
	}
}

func (container *iocContainer) RegisterTo(val interface{}, ifacePtr interface{}, lifecycle Lifecycle) {
	defer container.locker.Unlock()
	container.locker.Lock()

	typ := InterfaceOf(ifacePtr)
	if container.transient.Get(typ).IsValid() {
		lifecycle = Transient
	} else if container.singleton.Get(typ).IsValid() {
		lifecycle = Singleton
	}
	switch lifecycle {
	case Transient:
		container.transient.MapTo(val, ifacePtr)
		if DEBUG {
			consoleLog.Printf("Container.RegisterTo transient '%v'.\n", typ)
		}
	case Singleton:
		container.singleton.MapTo(val, ifacePtr)
		if DEBUG {
			consoleLog.Printf("Container.RegisterTo singleton '%v'.\n", typ)
		}
	}
}

func (container *iocContainer) Resolve(typ reflect.Type) reflect.Value {
	defer container.locker.RUnlock()
	container.locker.RLock()

	startTime := time.Now().UnixNano()
	var val = container.transient.Get(typ)
	if val.IsValid() {
		if initializer, ok := val.Interface().(Initializer); ok {
			container.Invoke(initializer.InitFunc())
		}
		endTime := time.Now().UnixNano()
		if DEBUG {
			consoleLog.Printf("Container.Resolve transient '%v' execute in %vms.\n", typ, float64(endTime-startTime)/float64(time.Millisecond))
		}
	} else if val = container.singleton.Get(typ); val.IsValid() {
		if initializer, ok := val.Interface().(Initializer); ok {
			container.Invoke(initializer.InitFunc())
		}
		endTime := time.Now().UnixNano()
		if DEBUG {
			consoleLog.Printf("Container.Resolve singleton '%v' execute in %vms.\n", typ, float64(endTime-startTime)/float64(time.Millisecond))
		}
	} else if container.parent != nil {
		val = container.parent.Resolve(typ)
	}
	return val
}

func (container *iocContainer) Invoke(f interface{}) ([]reflect.Value, error) {
	t := reflect.TypeOf(f)
	var in = make([]reflect.Value, t.NumIn()) //Panic if t is not kind of Func
	for i := 0; i < t.NumIn(); i++ {
		argType := t.In(i)
		val := container.Resolve(argType)
		if !val.IsValid() {
			return nil, fmt.Errorf("Value not found for type %v", argType)
		}
		in[i] = val
	}
	return reflect.ValueOf(f).Call(in), nil
}

func (container *iocContainer) SetParent(parent ReadonlyContainer) {
	container.parent = parent
}

// FromPtrTypeOf is to get real type from a pointer to value
func FromPtrTypeOf(obj interface{}) reflect.Type {
	realType := reflect.TypeOf(obj)
	for realType.Kind() == reflect.Ptr {
		realType = realType.Elem()
	}
	return realType
}

// FromPtrType is to get real type from a pointer to type
func FromPtrType(typ reflect.Type) reflect.Type {
	realType := typ
	for realType.Kind() == reflect.Ptr {
		realType = realType.Elem()
	}
	return realType
}

// InterfaceOf is to get interface from a pointer to value
func InterfaceOf(ifacePtr interface{}) reflect.Type {
	t := reflect.TypeOf(ifacePtr)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Interface {
		consoleLog.Panic("Called InterfaceOf with a value that is not a pointer to an interface. (*MyInterface)(nil)")
	}
	return t
}
