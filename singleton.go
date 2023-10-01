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
