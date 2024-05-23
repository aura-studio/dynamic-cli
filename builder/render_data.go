package builder

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/aura-studio/dynamic-cli/config"
)

type RenderData struct {
	Name        string
	Version     string
	Package     string
	FullPackage string
	Module      string
	House       string
	GoVersion   string
	NetRC       string
	LDFlags     string
}

func (r *RenderData) MustValid(renderData *RenderData) {
	for _, b := range []byte(r.Name) {
		if b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9' || b == '_' {
			continue
		}
		log.Panicf("invalid character '%s' in name", string(b))
	}
	for _, b := range []byte(r.Version) {
		if b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9' || b == '_' {
			continue
		}
		log.Panicf("invalid character '%s' in version", string(b))
	}
}

func NewRenderData(c config.Config) []*RenderData {
	ldFlags := ""
	if !c.Debug {
		ldFlags += "-s -w " // need an extra space
	}

	if len(c.Packages) == 0 {
		packagePath := c.Module
		packageName := packagePath[strings.LastIndex(packagePath, "/")+1:]
		name := strings.Join([]string{packageName, c.Commit}, "_")
		if len(c.Namespace) > 0 {
			name = strings.Join([]string{c.Namespace, name}, "_")
		}
		renderData := &RenderData{
			GoVersion:   c.GoVer,
			Name:        name,
			Version:     c.Commit,
			Package:     packageName,
			FullPackage: c.Module,
			Module:      c.Module,
			House:       fmt.Sprintf("%s/%s", runtime.Version(), c.WareHouse),
			NetRC:       c.NetRC,
			LDFlags:     ldFlags,
		}
		renderData.MustValid(renderData)
		return []*RenderData{renderData}
	}

	renderDatas := make([]*RenderData, len(c.Packages))
	for i, packagePath := range c.Packages {
		packageName := packagePath[strings.LastIndex(packagePath, "/")+1:]
		name := strings.Join([]string{packageName, c.Commit}, "_")
		if len(c.Namespace) > 0 {
			name = strings.Join([]string{c.Namespace, name}, "_")
		}
		renderData := &RenderData{
			GoVersion:   c.GoVer,
			Name:        name,
			Version:     c.Commit,
			Package:     packageName,
			FullPackage: strings.Join([]string{c.Module, packagePath}, "/"),
			Module:      c.Module,
			House:       fmt.Sprintf("%s/%s", runtime.Version(), c.WareHouse),
			NetRC:       c.NetRC,
			LDFlags:     ldFlags,
		}
		renderData.MustValid(renderData)
		renderDatas[i] = renderData
	}

	return renderDatas
}
