# ioc

Inversion of Control (IoC)

You can register a type as singleton or transient.
Also you can register a type mapping to an interface as singleton or transient.


## Usage

    go get github.com/berkaroad/ioc

## Change List

* 2016/8/31
    1. Remove readonly lock,
    2. Singleton instance's initialization called only once.
    3. Performance is 15% faster than last version.

* 2016/7/18
    1. first version.


## Performance

2 routine, 4 resolve action, 470,000 / sec,  compile at go 1.7

### Test code

    func main() {
        var requestContext = ioc.NewContainer()
        requestContext.SetParent(iocContainer)
        requestContext.RegisterTo(&productcategoryApp.ProductCategoryApplicationServiceImpl{}, (*application.ProductCategoryApplicationService)(nil), ioc.Transient)

        commandMQAdapter := new(provider.MyCommandMQProvider)
        processor := cqrs.NewCommandProcessor(commandMQAdapter)
        processor.RegisterMiddleware((*middleware.AuthCommandMiddleware)(nil))

        // execute count
        var exeCount = 1000000
        // concurrent routine
        var concurrentCount = 1
        for true {
            var wg *sync.WaitGroup = &sync.WaitGroup{}
            time.Sleep(300 * time.Millisecond)
            startTime := time.Now().UnixNano()
            for i := 0; i < concurrentCount; i++ {
                wg.Add(1)
                go func(wg1 *sync.WaitGroup) {
                    for j := 0; j < exeCount/concurrentCount; j++ {
                        requestContext.Invoke(func(productCategoryAppSvc application.ProductCategoryApplicationService, roContainer ioc.ReadonlyContainer) {
                            //processor.RegisterHandler(productCategoryAppSvc)
                        })
                    }
                    wg1.Done()
                }(wg)
            }
            wg.Wait()
            endTime := time.Now().UnixNano()
            consoleLog.Printf("[info] requestContext.Invoke for %d times with %d routines execute in %vms.\n", exeCount, concurrentCount, float64(endTime-startTime)/float64(time.Millisecond))
        }
    }

### Scenario 1:

1 routine, 3 times resolve singleton and 1 times resolve transient per code invoke, invoke 1,000,000 times.

Result:

    [commandprocessor] 2016/08/31 12:11:53 [info] requestContext.Invoke for 1000000 times with 1 routines execute in 3856.813759ms.
    [commandprocessor] 2016/08/31 12:11:57 [info] requestContext.Invoke for 1000000 times with 1 routines execute in 3973.131994ms.
    [commandprocessor] 2016/08/31 12:12:01 [info] requestContext.Invoke for 1000000 times with 1 routines execute in 3873.912007ms.
    [commandprocessor] 2016/08/31 12:12:06 [info] requestContext.Invoke for 1000000 times with 1 routines execute in 3940.871694ms.
    [commandprocessor] 2016/08/31 12:12:10 [info] requestContext.Invoke for 1000000 times with 1 routines execute in 3893.92014ms.
    [commandprocessor] 2016/08/31 12:12:15 [info] requestContext.Invoke for 1000000 times with 1 routines execute in 4853.996326ms.
    [commandprocessor] 2016/08/31 12:12:19 [info] requestContext.Invoke for 1000000 times with 1 routines execute in 4009.755085ms.
    [commandprocessor] 2016/08/31 12:12:24 [info] requestContext.Invoke for 1000000 times with 1 routines execute in 4077.67913ms.
    [commandprocessor] 2016/08/31 12:12:28 [info] requestContext.Invoke for 1000000 times with 1 routines execute in 3926.909983ms.
    [commandprocessor] 2016/08/31 12:12:32 [info] requestContext.Invoke for 1000000 times with 1 routines execute in 3899.237142ms.

### Scenario 2:

2 routine, 3 times resolve singleton and 1 times resolve transient per code invoke, invoke 1,000,000 times.

Result:

    [commandprocessor] 2016/08/31 12:14:36 [info] requestContext.Invoke for 1000000 times with 2 routines execute in 2034.126408ms.
    [commandprocessor] 2016/08/31 12:14:38 [info] requestContext.Invoke for 1000000 times with 2 routines execute in 2033.807996ms.
    [commandprocessor] 2016/08/31 12:14:40 [info] requestContext.Invoke for 1000000 times with 2 routines execute in 2051.833847ms.
    [commandprocessor] 2016/08/31 12:14:43 [info] requestContext.Invoke for 1000000 times with 2 routines execute in 2063.946131ms.
    [commandprocessor] 2016/08/31 12:14:45 [info] requestContext.Invoke for 1000000 times with 2 routines execute in 2032.77146ms.
    [commandprocessor] 2016/08/31 12:14:47 [info] requestContext.Invoke for 1000000 times with 2 routines execute in 2036.492861ms.
    [commandprocessor] 2016/08/31 12:14:50 [info] requestContext.Invoke for 1000000 times with 2 routines execute in 2091.975376ms.
    [commandprocessor] 2016/08/31 12:14:52 [info] requestContext.Invoke for 1000000 times with 2 routines execute in 2080.193925ms.
    [commandprocessor] 2016/08/31 12:14:55 [info] requestContext.Invoke for 1000000 times with 2 routines execute in 2030.570388ms.
    [commandprocessor] 2016/08/31 12:14:57 [info] requestContext.Invoke for 1000000 times with 2 routines execute in 2026.989203ms.

## Example

    package main

    import (
        "github.com/berkaroad/ioc"
    )

    type Interface1 interface {
        GetC2Name() string
    }

    type Interface2 interface {
        GetName() string
    }

    type Class1 struct {
        C2Name          string
        isInitialized   bool
    }

    func (c *Class1) InitFunc() interface{} {
        return func(c2 *Class2) {
            if !c.isInitialized {
                c.isInitialized = true
                c.C2Name = c2.Name
            }
        }
    }

    func (c *Class1) GetC2Name() string {
        return c.C2Name
    }

    type Class2 struct {
        Name            string
        isInitialized   bool
    }

    func (c *Class2) InitFunc() interface{} {
        return func() {
            if !c.isInitialized {
                c.isInitialized = true
                c.Name = "Tomcat"
            }
        }
    }

    func (c *Class2) GetName() string {
        return c.Name
    }

    func main() {
        var container = ioc.NewContainer()

        // Register class
        container.Register(&Class1{}, ioc.Singleton)
        container.Register(&Class2{Name: "Jerry Bai"}, ioc.Singleton)

        // Register class mapping to interface.
        container.RegisterTo(&Class1{}, (*Interface1)(nil), ioc.Transient)
        container.RegisterTo(&Class2{Name: "Jerry Bai"}, (*Interface2)(nil), ioc.Transient)

        // Like class's construction, inject class instance
        container.Invoke(func(c1 *Class1, c2 *Class2, roContainer ioc.ReadonlyContainer) {
            println("c1.C2Name=", c1.C2Name)
            println("c2.Name=", c2.Name)
        })

        // Like class's construction, inject interface instance
        container.Invoke(func(i1 Interface1, i2 Interface2, roContainer ioc.ReadonlyContainer) {
            println("i1.GetC2Name=()", i1.GetC2Name())
            println("i2.GetName=()", i2.GetName())
        })
    }
