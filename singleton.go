package ioc

import (
	"reflect"
)

type singletonContainer struct {
	valuemapper map[reflect.Type]reflect.Value
}

func (container *singletonContainer) Map(val interface{}) {
	container.valuemapper[reflect.TypeOf(val)] = reflect.ValueOf(val)
}

func (container *singletonContainer) MapTo(val interface{}, ifacePtr interface{}) {
	container.valuemapper[InterfaceOf(ifacePtr)] = reflect.ValueOf(val)
}

func (container *singletonContainer) Get(typ reflect.Type) reflect.Value {
	return container.valuemapper[typ]
}
