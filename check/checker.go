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
	// 检查操作系统的版本，如ubuntu22.04、alpine3.18等，不符合则警告
	// 约定：
	// - 若 TargetOS 为 go env GOOS 风格（linux/darwin/windows），则直接对比 GOOS
	// - 若 TargetOS 为发行版+版本（ubuntu22.04 / alpine3.18），则仅在 linux 下从 /etc/os-release 获取 ID+VERSION_ID
	// - 无法探测时仅告警，不中断
	if c.TargetOS == "" {
		return true
	}

	actualOS := strings.ToLower(strings.TrimSpace(GetOS()))
	if actualOS == "" {
		warnf("cannot detect GOOS (target=%s)", c.TargetOS)
		return false
	}

	// Generic GOOS match
	if isGenericGOOS(c.TargetOS) {
		expected := strings.ToLower(strings.TrimSpace(c.TargetOS))
		actualGOOS := actualOS
		// If getOS returns a distro descriptor (e.g. ubuntu22.04), treat it as linux.
		if !isGenericGOOS(actualGOOS) {
			actualGOOS = "linux"
		}
		if actualGOOS != expected {
			warnf("OS mismatch (target=%s actual=%s)", expected, actualGOOS)
			return false
		}
		okf("OS match (target=%s actual=%s)", expected, actualGOOS)
		return true
	}

	// Distro match only on linux
	if isGenericGOOS(actualOS) {
		// getOS couldn't detect distro; only got GOOS back.
		warnf("cannot detect distro version from /etc/os-release (target=%s)", c.TargetOS)
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
	// 检查Go的架构附带子版本，如arm64v8与arm64v7等，不符合则警告
	// 约定：
	// - TargetArch 支持：amd64/arm64/armv7/armv6/arm64v8 等
	// - 实际值优先来自 go env GOARCH/GOARM；必要时用 uname -m 辅助推断 arm64v8
	// - 允许“只写基础架构”的目标（如 arm64）匹配带后缀的实际（如 arm64v8）
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
	// Allow base match
	if baseArch(actualLower) == expected {
		okf("ARCH match (target=%s actual=%s)", expected, actualLower)
		return true
	}
	warnf("ARCH mismatch (target=%s actual=%s)", expected, actualLower)
	return false
}

func (c *Checker) checkCompiler() bool {
	// 检查Go的编译器版本，如go1.20.5等，不符合则警告
	// 约定：
	// - TargetCompiler 形如 go1.20.5 或 go1.20
	// - 若 TargetCompiler 是实际版本的前缀（go1.20 匹配 go1.20.5），则视为符合
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
	if strings.HasPrefix(actual, expected) {
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

func isGenericGOOS(osName string) bool {
	switch strings.ToLower(strings.TrimSpace(osName)) {
	case "linux", "darwin", "windows":
		return true
	default:
		return false
	}
}
func baseArch(arch string) string {
	arch = strings.ToLower(strings.TrimSpace(arch))
	// Examples: arm64v8 -> arm64, armv7 -> arm
	if strings.HasPrefix(arch, "arm64") {
		return "arm64"
	}
	if strings.HasPrefix(arch, "armv") {
		return "arm"
	}
	return arch
}
