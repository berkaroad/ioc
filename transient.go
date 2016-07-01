package ioc

import (
	"reflect"
)

type transientContainer struct {
	typemapper map[reflect.Type]reflect.Type
}

func (container *transientContainer) Map(val interface{}) {
	typ := reflect.TypeOf(val)
	container.typemapper[typ] = FromPtrType(typ)
}

func (container *transientContainer) MapTo(val interface{}, ifacePtr interface{}) {
	typ := InterfaceOf(ifacePtr)
	container.typemapper[typ] = FromPtrTypeOf(val)
}

func (container *transientContainer) Get(typ reflect.Type) reflect.Value {
	var val reflect.Value
	if realTyp := container.typemapper[typ]; realTyp != nil {
		val = reflect.New(realTyp)
		if typ.Kind() != reflect.Ptr && typ.Kind() != reflect.Interface {
			val = val.Elem()
		}
	}
	return val
}
