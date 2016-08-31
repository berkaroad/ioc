package ioc

import (
	"reflect"
	"sync"
)

type transientContainer struct {
	locker     *sync.RWMutex
	typemapper map[reflect.Type]reflect.Type
}

func newTransientContainer() *transientContainer {
	c := new(transientContainer)
	c.locker = &sync.RWMutex{}
	c.typemapper = make(map[reflect.Type]reflect.Type)
	return c
}

func (container *transientContainer) Map(val interface{}) {
	defer container.locker.Unlock()
	container.locker.Lock()

	typ := reflect.TypeOf(val)
	container.typemapper[typ] = FromPtrType(typ)
}

func (container *transientContainer) MapTo(val interface{}, ifacePtr interface{}) {
	defer container.locker.Unlock()
	container.locker.Lock()

	typ := InterfaceOf(ifacePtr)
	container.typemapper[typ] = FromPtrTypeOf(val)
}

func (container *transientContainer) Get(typ reflect.Type, invoker func(f interface{})) reflect.Value {
	// defer container.locker.RUnlock()
	// container.locker.RLock()

	if realTyp := container.typemapper[typ]; realTyp != nil {
		val := reflect.New(realTyp)
		if typ.Kind() != reflect.Ptr && typ.Kind() != reflect.Interface {
			val = val.Elem()
		}

		if invoker != nil {
			if initializer, ok := val.Interface().(Initializer); ok {
				invoker(initializer.InitFunc())
			}
		}
		return val
	}
	return reflect.Value{}
}

func (container *transientContainer) Exists(typ reflect.Type) bool {
	// defer container.locker.RUnlock()
	// container.locker.RLock()

	return container.typemapper[typ] != nil
}
