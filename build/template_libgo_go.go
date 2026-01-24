package build

func init() {
	templateMap["{{.Dir}}/libgo_{{.Name}}/libgo.go"] = templateLibgoGo
}

const templateLibgoGo = `package main
/*
#cgo CFLAGS: -I../
#cgo LDFLAGS: -L../ -lcgo_{{.Name}}
#include "{{.Dir}}/libcgo_{{.Name}}.h"
#include "stdlib.h"
*/
import "C"
import (
	"encoding/json"
	"unsafe"
)

// Meta constants injected at build time via -ldflags -X
var (
	MetaModulePath string
	MetaCommitID   string
	MetaBuildTs    string
)

type tunnel struct{}

func (t tunnel) Meta() string {
	meta := map[string]string{
		"module_path": MetaModulePath,
		"commit_id":   MetaCommitID,
		"build_ts":    MetaBuildTs,
	}
	data, _ := json.Marshal(meta)
	return string(data)
}

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
