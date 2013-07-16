scgen
=====

Regexp based generator of gorilla/rpc-compatible service/client code for a specified object.

Mission
=====

scgen simplifies gorilla/rpc-compatible service (and client for it) generation for the specified funcs of a specified type.

Usage example
=====

Provided you have the following file 'example.txt':

```go
// ...
func (t *T) Func1(p1 Param1) (*SomeStruct, error) {
    ...
}

// aaaaaaaaaa
func (t *T) SomeFunc2(p2 *Param2, p3 Param3) (error) {
    a
    b
    c
}

...
func (t *T) Func3(p3 string, p4 interface{}) (bool, error) {
    ......
}

// ...
func (t *T) SomeFunc4() ([]string, error) {
    ...
}
```

you can simplify generation of your service by calling

```
./scgen -fsign="example.txt" -obj="t" -objt="T"
```

This will generate gorilla/rpc-compatible service func declarations:

```go
//------------------------------------------------
// ▢ Func1
//------------------------------------------------

type Func1_Args struct { 
    P1 Param1
}

type Func1_Result struct { 
    Par0 *SomeStruct
}

func (service *Service) Func1(r *http.Request, args *Func1_Args, reply *Func1_Result) (error) {
    logger.Debugf("P1: '%#v'",
                  args.P1)
    
    r0, err := service.T.Func1(args.P1)

    if err != nil { logger.Error(err); return err}
    reply.Par0 = r0

    return nil
}

//------------------------------------------------
// ▢ SomeFunc2
//------------------------------------------------

type SomeFunc2_Args struct { 
    P2 *Param2
    P3 Param3
}

type SomeFunc2_Result struct { 
}

func (service *Service) SomeFunc2(r *http.Request, args *SomeFunc2_Args, reply *SomeFunc2_Result) (error) {
    logger.Debugf("P2: '%#v', P3: '%#v'",
                  args.P2, args.P3)
    
    err := service.T.SomeFunc2(args.P2, args.P3)

    if err != nil { logger.Error(err); return err}

    return nil
}

//------------------------------------------------
// ▢ Func3
//------------------------------------------------

type Func3_Args struct { 
    P3 string
    P4 interface{}
}

type Func3_Result struct { 
    Par0 bool
}

func (service *Service) Func3(r *http.Request, args *Func3_Args, reply *Func3_Result) (error) {
    logger.Debugf("P3: '%s', P4: '%#v'",
                  args.P3, args.P4)
    
    r0, err := service.T.Func3(args.P3, args.P4)

    if err != nil { logger.Error(err); return err}
    reply.Par0 = r0

    return nil
}

//------------------------------------------------
// ▢ SomeFunc4
//------------------------------------------------

type SomeFunc4_Args struct { 
}

type SomeFunc4_Result struct { 
    Par0 []string
}

func (service *Service) SomeFunc4(r *http.Request, args *SomeFunc4_Args, reply *SomeFunc4_Result) (error) {
    logger.Debugf("")
    
    r0, err := service.T.SomeFunc4()

    if err != nil { logger.Error(err); return err}
    reply.Par0 = r0

    return nil
}

```

and client declarations:

```go
//------------------------------------------------
// ▢ Func1
//------------------------------------------------

func (client *ServiceClient) Func1(p1 Param1) (*SomeStruct, error) {
    args := service.Func1_Args{ p1 }
    var r service.Func1_Result

    e := client.GetResult(ServiceName + "Func1", &args, &r)
    return r.Par0, e
}


//------------------------------------------------
// ▢ SomeFunc2
//------------------------------------------------

func (client *ServiceClient) SomeFunc2(p2 *Param2, p3 Param3) (error) {
    args := service.SomeFunc2_Args{ p2, p3 }
    var r service.SomeFunc2_Result

    e := client.GetResult(ServiceName + "SomeFunc2", &args, &r)
    return e
}


//------------------------------------------------
// ▢ Func3
//------------------------------------------------

func (client *ServiceClient) Func3(p3 string, p4 interface{}) (bool, error) {
    args := service.Func3_Args{ p3, p4 }
    var r service.Func3_Result

    e := client.GetResult(ServiceName + "Func3", &args, &r)
    return r.Par0, e
}


//------------------------------------------------
// ▢ SomeFunc4
//------------------------------------------------

func (client *ServiceClient) SomeFunc4() ([]string, error) {
    args := service.SomeFunc4_Args{  }
    var r service.SomeFunc4_Result

    e := client.GetResult(ServiceName + "SomeFunc4", &args, &r)
    return r.Par0, e
}
```

After that all you need to do manually is to create the Service/ServiceClient types, their constructors, and correct imports, needed for all your types.
