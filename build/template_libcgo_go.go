package build

func init() {
	templateMap["{{.Dir}}/libcgo_{{.Name}}/libcgo.go"] = templateLibcgoGo
}

const templateLibcgoGo = `package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"unsafe"

	src "{{if eq .Package "."}}{{.Module}}{{else}}{{.Module}}/{{.Package}}{{end}}"
)

// Meta constants injected at build time via -ldflags -X
var (
	MetaSourceModule   string
	MetaSourceVersion  string
	MetaSourceBuilt    string
	MetaToolchainOS       string
	MetaToolchainArch     string
	MetaToolchainCompiler string
	MetaToolchainVariant  string
)

var tunnel = src.Tunnel

//export dynamic_cgo_{{.Name}}_meta
func dynamic_cgo_{{.Name}}_meta() *C.char {
	meta := map[string]interface{}{
		"source": map[string]string{
			"module":  MetaSourceModule,
			"version": MetaSourceVersion,
			"built":   MetaSourceBuilt,
		},
		"toolchain": map[string]string{
			"os":       MetaToolchainOS,
			"arch":     MetaToolchainArch,
			"compiler": MetaToolchainCompiler,
			"variant":  MetaToolchainVariant,
		},
	}
	data, _ := json.MarshalIndent(meta, "", "  ")
	return C.CString(string(data))
}

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
	tunnel.Close()
}

func main() {}
`
