package build

func init() {
	templateMap["{{.Dir}}/libcgo_{{.Name}}/libcgo.go"] = `package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	src "{{.Module}}"
)

var tunnel = src.Tunnel

//export dynamic_cgo_{{.Name}}_init
func dynamic_cgo_{{.Name}}_init() {
	tunnel.Init()
}

//export dynamic_cgo_{{.Name}}_invoke
func dynamic_cgo_{{.Name}}_invoke(route_cstr *C.char, req_cstr *C.char) *C.char {
	route := C.GoString(route_cstr)
	C.free(unsafe.Pointer(route_cstr))

	req := C.GoString(req_cstr)
	C.free(unsafe.Pointer(req_cstr))

	return C.CString(tunnel.Invoke(route, req))
}

//export dynamic_cgo_{{.Name}}_close
func dynamic_cgo_{{.Name}}_close() {
	tu
}
