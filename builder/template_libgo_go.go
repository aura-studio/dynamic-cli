package builder

func init() {
	templateMap["{{.House}}/{{.Name}}/libgo_{{.Name}}/libgo.go"] = `package main

/*
#cgo CFLAGS: -I{{.House}}/{{.Name}}/
#cgo LDFLAGS: -L{{.House}}/{{.Name}}/ -lcgo_{{.Name}}
#include "{{.House}}/{{.Name}}/libcgo_{{.Name}}.h"
#include "stdlib.h"
*/
import "C"
import "unsafe"

type tunnel struct{}

func (t tunnel) Init() {
	C.dynamic_cgo_{{.Name}}_init()
}

func (t tunnel) Invoke(route string, req string) string {
	rsp_cstr := C.dynamic_cgo_{{.Name}}_invoke(C.CString(route), C.CString(req))
	rsp := C.GoString(rsp_cstr)
	C.free(unsafe.Pointer(rsp_cstr))
	return rsp
}

func (t tunnel) Close() {
	C.dynamic_cgo_{{.Name}}_close()
}

var Tunnel tunnel

func main() {}
`
}
