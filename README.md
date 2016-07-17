# ioc
Inversion of Control (IoC)

You can register a type as singleton or transient.
Also you can register a type mapping to an interface as singleton or transient.


## Usage

    go get github.com/berkaroad/ioc

## Performance

2 routine, 4 resolve action, 400,000 / sec

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

### Scenario 1: 1 routine, 3 times resolve singleton and 1 times resolve transient per code invoke.

Result:

    [commandprocessor] 2016/07/17 11:31:29 [info] requestContext.Invoke for 1000000 times with 1 routines execute in 4971.1971ms.
    [commandprocessor] 2016/07/17 11:31:34 [info] requestContext.Invoke for 1000000 times with 1 routines execute in 4951.494214ms.
    [commandprocessor] 2016/07/17 11:31:39 [info] requestContext.Invoke for 1000000 times with 1 routines execute in 4954.376794ms.
    
### Scenario 2: 2 routine, 3 times resolve singleton and 1 times resolve transient per code invoke.

    [commandprocessor] 2016/07/17 11:23:50 [info] requestContext.Invoke for 1000000 times with 2 routines execute in 2779.720723ms.
    [commandprocessor] 2016/07/17 11:23:53 [info] requestContext.Invoke for 1000000 times with 2 routines execute in 2719.810844ms.
    [commandprocessor] 2016/07/17 11:23:56 [info] requestContext.Invoke for 1000000 times with 2 routines execute in 2734.028326ms.

Result:



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