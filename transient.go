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
