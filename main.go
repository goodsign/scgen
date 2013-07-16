package main 

import (
    "os"
    "fmt"
    "flag"
    "io/ioutil"
)
const (
    HelpFile = "example.txt"
)

// Flags
var (
    f     = flag.String("fsign",  "example.txt",      "File containing original object and method signatures")
    objT  = flag.String("objt",   "",                  "Original object type name")
    obj   = flag.String("obj",    "",                  "Original object name")
    sname = flag.String("sname",  "Service",           "Name of the service")
    cname = flag.String("cname",  "ServiceClient",     "Name of the client")
)

func main() {
    flag.Parse()

    if len(*f) == 0 || len(*sname) == 0 {
        fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
        flag.PrintDefaults()

        fmt.Println("Example of 'fsign' file contents:")
        help, err := ioutil.ReadFile(HelpFile)
        if err != nil {
            panic(err)
        }
        fmt.Println(string(help))
        os.Exit(0)
    }

    str, err := ioutil.ReadFile(*f)
    if err != nil {
        fmt.Println("ERROR: no such file:", *f)
        os.Exit(-1)
    }

    fmt.Print(TransformText(string(str), *obj, *objT, *sname, *cname))
}