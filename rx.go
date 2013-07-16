package main

import (
    "fmt"
    "bytes"
    "regexp"
    "strings"
    "io/ioutil"
    "text/template"
)

const (
    TemplateIncompatibilityMsg = "Template is not consistent with data:"
)

/*
==================================================

HOW IT WORKS:

For each func, that is found using a funcRx regexp: its name, params, result args are extracted.
After that these params are wrapped into special DTOs and passed to the text templates, stored
in the same directory as the code (with *.tmpl extension).

After that the final output is created by concatenating results of these templates output for
each func.
==================================================
*/

type rxWrapper struct {
    *regexp.Regexp
}

// The RX which is used to find funcs
var funcRx = rxWrapper{regexp.MustCompile(`\s*?func\s*?\(\s*?(?P<rname>\w+?)\s+?\*?(?P<rtype>[\w]+?)\s*?\)\s+(?P<fname>\w+)\s*?\(\s*?(?P<fparams>[\s\S]*?)\)\s*?\(?\s*?(?P<freturn>[\s\S]*?)\)?\s*?\{`)}

// RX to find params and result args
var parRx = rxWrapper{regexp.MustCompile(`(?P<pname>\w+\s+)?(?P<ptype>[\w\*\{\}\[\]\.]+)`)}

var argsTmp *template.Template
var retvalTmp *template.Template
var servFuncTmp *template.Template
var clientFuncTmp *template.Template
var servTmp *template.Template
var clientTmp *template.Template

func loadTemplate(name string) *template.Template {
    argsB, err := ioutil.ReadFile(name + ".tmpl")
    if err != nil { panic (err) }
    return template.Must(template.New("name").Parse(string(argsB)))
}

func init() {
    // Load templates
    argsTmp = loadTemplate("args")
    retvalTmp = loadTemplate("retval")
    servFuncTmp = loadTemplate("serv_func")
    clientFuncTmp = loadTemplate("client_func")
    servTmp = loadTemplate("serv")
    clientTmp = loadTemplate("client")
}


//------------------------------------------------
// ▢ Utilities
//------------------------------------------------

// FindStringSubmatchMap is a utility func wrapper for regexp.FindStringSubmatch
// but outputs result as a map with keys = capture names, values = capture values.
func (r *rxWrapper) FindStringSubmatchMap(s string) map[string]string {
    asm := r.FindAllStringSubmatchMap(s)

    if len(asm) == 0 {
        return nil
    }

    return asm[0]
}

// FindAllStringSubmatchMap is a utility func wrapper for regexp.FindAllStringSubmatch
// but outputs results as maps with keys = capture names, values = capture values.
func (r *rxWrapper) FindAllStringSubmatchMap(s string) []map[string]string {
    matches := r.FindAllStringSubmatch(s, -1)
    captures := make([]map[string]string, len(matches))

    for i, match := range matches {
        captures[i] = make(map[string]string)
        for j, name := range r.SubexpNames() {
            if j == 0 {
                continue
            }
            captures[i][name] = match[j]
        }
    }
    return captures
}

//------------------------------------------------
// ▢ Templates and rendering
//------------------------------------------------

// ArgsRenderData is a DTO for "./args.tmpl"
type ArgsRenderData struct {
    Args    []ParamDesc
    FName   string
}

// generateArgsString executes "./args.tmpl" on a DTO with given func names and param descriptions.
func generateArgsString(fName string, pd []ParamDesc) string {
    adata := ArgsRenderData{pd, fName}

    var b bytes.Buffer
    err := argsTmp.Execute(&b, adata)

    if err != nil {
        panic(TemplateIncompatibilityMsg + err.Error())
    }

    return b.String()
}

// ResultsRenderData is a DTO for "./result.tmpl"
type ResultsRenderData struct {
    RetValues   []ParamDesc
    FName       string
}

// generateRetValsString executes "./retval.tmpl" on a DTO with given func names and param descriptions.
func generateRetValsString(fName string, rd []ParamDesc) string {
    adata := ResultsRenderData{rd, fName}

    var b bytes.Buffer
    err := retvalTmp.Execute(&b, adata)

    if err != nil {
        panic(TemplateIncompatibilityMsg + err.Error())
    }

    return b.String()
}

// ServiceFuncData is a DTO for "./serv_func.tmpl"
type ServiceFuncData struct {
    ServName    string
    ObjName     string
    ObjTName    string
    FName       string
    Args        []ParamDesc
    RetValues   []ParamDesc     
}

// generateServiceFuncString executes "./serv_func.tmpl" on a DTO with given func names and param descriptions.
func generateServiceFuncString(fName string, args []ParamDesc, retvals []ParamDesc,
                               servName string, objName string, objTName string) string {
    adata := ServiceFuncData{servName, objName, objTName, fName, args, retvals}
    var b bytes.Buffer
    err := servFuncTmp.Execute(&b, adata)

    if err != nil {
        panic(TemplateIncompatibilityMsg + err.Error())
    }

    return b.String()
}

// ServiceData is a DTO for "./serv.tmpl"
type ServiceData struct {
    ServName    string
    FuncArgs    []string 
    FuncResults []string
    Funcs       []string
}

// generateServiceString executes "./service.tmpl" on a DTO with given func names and param descriptions.
func generateServiceString(servName string, funcArgs []string, funcRes []string, funcs []string) string {
    adata := ServiceData{servName, funcArgs, funcRes, funcs}
    var b bytes.Buffer
    err := servTmp.Execute(&b, adata)

    if err != nil {
        panic(TemplateIncompatibilityMsg + err.Error())
    }

    return b.String()
}

// ClientFuncData is a DTO for "./client_func.tmpl"
type ClientFuncData struct {
    ServName    string
    ObjName     string
    ObjTName    string
    ClientName  string
    FName       string
    Args        []ParamDesc
    RetValues   []ParamDesc     
}

// generateClientFuncString executes "./serv_func.tmpl" on a DTO with given func names and param descriptions.
func generateClientFuncString(fName string, args []ParamDesc, retvals []ParamDesc,
                              servName string, objName string, objTName string, clientName string) string {
    adata := ClientFuncData{servName, objName, objTName, clientName, fName, args, retvals}
    var b bytes.Buffer
    err := clientFuncTmp.Execute(&b, adata)

    if err != nil {
        panic(TemplateIncompatibilityMsg + err.Error())
    }

    return b.String()
}

// ClientData is a DTO for "./serv.tmpl"
type ClientData struct {
    ServName    string
    ClientFuncs []string
}

// generateClientString executes "./client.tmpl" on a DTO with given func names and param descriptions.
func generateClientString(servName string, clientFuncs []string) string {
    adata := ClientData{servName, clientFuncs}
    var b bytes.Buffer
    err := clientTmp.Execute(&b, adata)

    if err != nil {
        panic(TemplateIncompatibilityMsg + err.Error())
    }

    return b.String()
}

//------------------------------------------------
// ▢ Extraction
//------------------------------------------------

// ParamDesc describes a parameter/returned value in a simplest way
type ParamDesc struct {
    Pname string 
    Ptype string
}

func (pd *ParamDesc) GetCapPname() string {
    return strings.ToUpper(pd.Pname[:1]) + pd.Pname[1:]
}

// GetLoggerFormatId returns format id for current type
func (pd *ParamDesc) GetLoggerFormatId() string {
    switch(pd.  Ptype) {
    case "int", "uint", "int32", "uint32", "int64", "uint64", "int8", "uint8":
        return "%d"
    case "float32", "float64":
        return "%f"
    case "string":
        return "%s"
    }

    return "%#v"
}

// "a param1, B *Param2, c interface{}" -> [{a param1}, {B *Param2}, {c interface{}}]
func extractParams(paramStr string) []ParamDesc {
    ps := parRx.FindAllStringSubmatchMap(paramStr)

    pds := make([]ParamDesc, len(ps))
    for i, p := range ps {
        pds[i] = ParamDesc{strings.TrimSpace(p["pname"]), strings.TrimSpace(p["ptype"])}
    }

    return pds
}

// processFunc does all the work for each func. 
func processFunc(fName string, argsStr string, retValsStr string, objName string, objTName string,
                 servName string, clientName string) (args string, rets string, fnc string, clientfnc string){
    retVals := extractParams(retValsStr)

    // remove last 'error' retval
    retVals = retVals[:len(retVals) - 1]

    // give names to unnamed retvals
    for i, v := range retVals {
        if len(strings.TrimSpace(v.Pname)) == 0 {
            retVals[i].Pname = fmt.Sprintf("par%d", i)
        }
    }

    argVals := extractParams(argsStr)

    return generateArgsString(fName, argVals), 
           generateRetValsString(fName, retVals), 
           generateServiceFuncString(fName, argVals, retVals, servName, objName, objTName),
           generateClientFuncString(fName, argVals, retVals, servName, objName, objTName, clientName)
}

// transformText transforms original code file containing correct list of methods
// of some object (objName) to the same number of service/client methods
// for a service with name 'servName' and client with name 'clientName'
func TransformText(input string, objName string, objTName string, servName string, clientName string) string {
    asm := funcRx.FindAllStringSubmatchMap(input)

    var sargs    []string
    var sretvals []string
    var sfuncs   []string
    var sclfuncs []string

    for _, funcL := range asm {
        if (funcL["rname"] == objName || len(objName) == 0) &&
           (funcL["rtype"] == objTName || len(objTName) == 0) {
            args, retvals, funcs, clFuncs := processFunc(funcL["fname"], funcL["fparams"], funcL["freturn"],
                                                         funcL["rname"], funcL["rtype"], servName, clientName)
            sargs = append(sargs, args)
            sretvals = append(sretvals, retvals)
            sfuncs = append(sfuncs, funcs)
            sclfuncs = append(sclfuncs, clFuncs)
        }
    }

    fmt.Println(generateServiceString(servName, sargs, sretvals, sfuncs))
    fmt.Println(generateClientString(servName, sclfuncs))

    return ""
}