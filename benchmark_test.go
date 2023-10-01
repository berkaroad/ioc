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
	"context"
	"testing"
)

func BenchmarkInjectToFunc(b *testing.B) {
	reset()
	AddSingleton[ProductCategoryRepository](&ProductCategoryRepositoryImpl{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc := &ProductCategoryApplicationServiceImpl{}
		Inject(svc.Initialize)
		svc.Get(context.TODO(), "123")
	}
}

func BenchmarkInjectToStruct(b *testing.B) {
	reset()
	AddSingleton[ProductCategoryRepository](&ProductCategoryRepositoryImpl{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc := &ProductCategoryApplicationServiceImpl{}
		Inject(svc)
		svc.Get(context.TODO(), "123")
	}
}

type ProductCategoryApplicationService interface {
	Get(ctx context.Context, id string) ProductCategory
}

var _ ProductCategoryApplicationService = (*ProductCategoryApplicationServiceImpl)(nil)

type ProductCategoryApplicationServiceImpl struct {
	Resolver Resolver
	Repo     ProductCategoryRepository `ioc-inject:"true"`
}

func (svc *ProductCategoryApplicationServiceImpl) Initialize(resolver Resolver, repo ProductCategoryRepository) {
	svc.Resolver = resolver
	svc.Repo = repo
}

func (svc *ProductCategoryApplicationServiceImpl) Get(ctx context.Context, id string) ProductCategory {
	return svc.Repo.Get(id)
}

type ProductCategoryRepository interface {
	Get(id string) ProductCategory
}

var _ ProductCategoryRepository = (*ProductCategoryRepositoryImpl)(nil)

type ProductCategoryRepositoryImpl struct{}

func (repo *ProductCategoryRepositoryImpl) Get(id string) ProductCategory {
	return ProductCategory{
		ID: id,
	}
}

type ProductCategory struct {
	ID   string
	Name string
}
