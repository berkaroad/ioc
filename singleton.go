package ioc

import (
	"reflect"
	"sync"
)

type stateValue struct {
	isInitialized bool
	val           reflect.Value
}

type singletonContainer struct {
	locker      *sync.RWMutex
	valuemapper map[reflect.Type]*stateValue
}

func newSingletonContainer() *singletonContainer {
	c := new(singletonContainer)
	c.locker = &sync.RWMutex{}
	c.valuemapper = make(map[reflect.Type]*stateValue)
	return c
}

func (container *singletonContainer) Map(val interface{}) {
	defer container.locker.Unlock()
	container.locker.Lock()

	container.valuemapper[reflect.TypeOf(val)] = &stateValue{val: reflect.ValueOf(val)}
}

func (container *singletonContainer) MapTo(val interface{}, ifacePtr interface{}) {
	defer container.locker.Unlock()
	container.locker.Lock()

	container.valuemapper[InterfaceOf(ifacePtr)] = &stateValue{val: reflect.ValueOf(val)}
}

func (container *singletonContainer) Get(typ reflect.Type, invoker func(f interface{})) reflect.Value {
	// defer container.locker.RUnlock()
	// container.locker.RLock()

	if stateVal := container.valuemapper[typ]; stateVal != nil {
		if !stateVal.isInitialized && invoker != nil {
			if initializer, ok := stateVal.val.Interface().(Initializer); ok {
				invoker(initializer.InitFunc())
				stateVal.isInitialized = true
			}
		}
		return stateVal.val
	}
	return reflect.Value{}
}

func (container *singletonContainer) Exists(typ reflect.Type) bool {
	// defer container.locker.RUnlock()
	// container.locker.RLock()

	return container.valuemapper[typ] != nil
}
