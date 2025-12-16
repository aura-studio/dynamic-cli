package check

import (
	"fmt"
	"os"
	"strings"

	"github.com/aura-studio/dynamic-cli/config"
)

type Checker struct {
	TargetOS       string
	TargetArch     string
	TargetCompiler string
}

func NewChecker(proc config.Procedure) *Checker {
	return &Checker{
		TargetOS:       proc.Toolchain.OS,
		TargetArch:     proc.Toolchain.Arch,
		TargetCompiler: proc.Toolchain.Compiler,
	}
}

func (c *Checker) Run() bool {
	okOS := c.checkOS()
	okArch := c.checkArch()
	okCompiler := c.checkCompiler()
	if okOS && okArch && okCompiler {
		return true
	}
	return false
}

func (c *Checker) checkOS() bool {
	// 直接比对 env.GetOS() 的返回值（严格相等）。
	if c.TargetOS == "" {
		return true
	}

	actualOS := strings.ToLower(strings.TrimSpace(GetOS()))
	if actualOS == "" {
		warnf("cannot detect GOOS (target=%s)", c.TargetOS)
		return false
	}

	expected := strings.ToLower(strings.TrimSpace(c.TargetOS))
	if actualOS != expected {
		warnf("OS mismatch (target=%s actual=%s)", expected, actualOS)
		return false
	}
	okf("OS match (target=%s actual=%s)", expected, actualOS)
	return true
}

func (c *Checker) checkArch() bool {
	// 直接比对 env.GetArch() 的返回值（严格相等）。
	if c.TargetArch == "" {
		return true
	}

	actual := GetArch()
	if actual == "" {
		warnf("cannot detect arch (target=%s)", c.TargetArch)
		return false
	}

	expected := strings.ToLower(strings.TrimSpace(c.TargetArch))
	actualLower := strings.ToLower(actual)

	if expected == actualLower {
		okf("ARCH match (target=%s actual=%s)", expected, actualLower)
		return true
	}
	warnf("ARCH mismatch (target=%s actual=%s)", expected, actualLower)
	return false
}

func (c *Checker) checkCompiler() bool {
	// 检查Go的编译器版本，如go1.20.5等，不符合则警告（严格相等）。
	if c.TargetCompiler == "" {
		return true
	}
	actual := GetComplier()
	if actual == "" {
		warnf("cannot detect GOVERSION (target=%s)", c.TargetCompiler)
		return false
	}
	expected := strings.TrimSpace(c.TargetCompiler)
	if actual == expected {
		okf("COMPILER match (target=%s actual=%s)", expected, actual)
		return true
	}
	warnf("COMPILER mismatch (target=%s actual=%s)", expected, actual)
	return false
}

func warnf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "fail: "+format+"\n", args...)
}

func okf(format string, args ...any) {
	fmt.Printf("pass: "+format+"\n", args...)
}
