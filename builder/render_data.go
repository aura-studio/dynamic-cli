package builder

import (
	"log"
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
		if b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '_' {
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
		pkg := c.Module[strings.LastIndex(c.Module, "/")+1:]
		name := pkg
		if len(c.Namespace) > 0 {
			name = strings.Join([]string{c.Namespace, pkg}, "_")
		}
		renderData := &RenderData{
			GoVersion:   c.GoVer,
			Name:        name,
			Version:     c.Commit,
			Package:     pkg,
			FullPackage: c.Module,
			Module:      c.Module,
			House:       c.WareHouse,
			NetRC:       c.NetRC,
			LDFlags:     ldFlags,
		}
		renderData.MustValid(renderData)
		return []*RenderData{renderData}
	}

	renderDatas := make([]*RenderData, len(c.Packages))
	for i, pkg := range c.Packages {
		name := pkg[strings.LastIndex(pkg, "/")+1:]
		if len(c.Namespace) > 0 {
			name = strings.Join([]string{c.Namespace, name}, "_")
		}
		renderData := &RenderData{
			GoVersion:   c.GoVer,
			Name:        name,
			Version:     c.Commit,
			Package:     pkg,
			FullPackage: strings.Join([]string{c.Module, pkg}, "/"),
			Module:      c.Module,
			House:       c.WareHouse,
			NetRC:       c.NetRC,
			LDFlags:     ldFlags,
		}
		renderData.MustValid(renderData)
		renderDatas[i] = renderData
	}

	return renderDatas
}
