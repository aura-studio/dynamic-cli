package check

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

func (c *Checker) Run() {
	c.checkOS()
	c.checkArch()
	c.checkCompiler()
}

func (c *Checker) checkOS() {
	// 检查操作系统的版本，如ubuntu22.04、alpine3.18等，不符合则警告
	// 约定：
	// - 若 TargetOS 为 go env GOOS 风格（linux/darwin/windows），则直接对比 GOOS
	// - 若 TargetOS 为发行版+版本（ubuntu22.04 / alpine3.18），则仅在 linux 下从 /etc/os-release 获取 ID+VERSION_ID
	// - 无法探测时仅告警，不中断
	if c.TargetOS == "" {
		return
	}

	actualGOOS := goEnv("GOOS")
	if actualGOOS == "" {
		warnf("toolchain cannot detect GOOS (target=%s)", c.TargetOS)
		return
	}

	// Generic GOOS match
	if isGenericGOOS(c.TargetOS) {
		if actualGOOS != c.TargetOS {
			warnf("toolchain OS mismatch (target=%s actual=%s)", c.TargetOS, actualGOOS)
			return
		}
		okf("toolchain OS match (target=%s actual=%s)", c.TargetOS, actualGOOS)
		return
	}

	// Distro match only on linux
	if actualGOOS != "linux" {
		warnf("toolchain OS mismatch (target=%s actual=%s)", c.TargetOS, actualGOOS)
		return
	}

	id, ver := readOSRelease()
	if id == "" || ver == "" {
		warnf("toolchain cannot detect distro version from /etc/os-release (target=%s)", c.TargetOS)
		return
	}
	actual := strings.ToLower(id) + strings.TrimSpace(ver)
	expected := strings.ToLower(strings.TrimSpace(c.TargetOS))
	if actual != expected {
		warnf("toolchain OS mismatch (target=%s actual=%s)", expected, actual)
		return
	}
	okf("toolchain OS match (target=%s actual=%s)", expected, actual)
}

func (c *Checker) checkArch() {
	// 检查Go的架构附带子版本，如arm64v8与arm64v7等，不符合则警告
	// 约定：
	// - TargetArch 支持：amd64/arm64/armv7/armv6/arm64v8 等
	// - 实际值优先来自 go env GOARCH/GOARM；必要时用 uname -m 辅助推断 arm64v8
	// - 允许“只写基础架构”的目标（如 arm64）匹配带后缀的实际（如 arm64v8）
	if c.TargetArch == "" {
		return
	}

	actual := detectArchWithVariant()
	if actual == "" {
		warnf("toolchain cannot detect arch (target=%s)", c.TargetArch)
		return
	}

	expected := strings.ToLower(strings.TrimSpace(c.TargetArch))
	actualLower := strings.ToLower(actual)

	if expected == actualLower {
		okf("toolchain ARCH match (target=%s actual=%s)", expected, actualLower)
		return
	}
	// Allow base match
	if baseArch(actualLower) == expected {
		okf("toolchain ARCH match (target=%s actual=%s)", expected, actualLower)
		return
	}
	warnf("toolchain ARCH mismatch (target=%s actual=%s)", expected, actualLower)
}

func (c *Checker) checkCompiler() {
	// 检查Go的编译器版本，如go1.20.5等，不符合则警告
	// 约定：
	// - TargetCompiler 形如 go1.20.5 或 go1.20
	// - 若 TargetCompiler 是实际版本的前缀（go1.20 匹配 go1.20.5），则视为符合
	if c.TargetCompiler == "" {
		return
	}
	actual := goEnv("GOVERSION")
	if actual == "" {
		warnf("toolchain cannot detect GOVERSION (target=%s)", c.TargetCompiler)
		return
	}
	expected := strings.TrimSpace(c.TargetCompiler)
	if actual == expected {
		okf("toolchain COMPILER match (target=%s actual=%s)", expected, actual)
		return
	}
	if strings.HasPrefix(actual, expected) {
		okf("toolchain COMPILER match (target=%s actual=%s)", expected, actual)
		return
	}
	warnf("toolchain COMPILER mismatch (target=%s actual=%s)", expected, actual)
}

func warnf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "fail: "+format+"\n", args...)
}

func okf(format string, args ...any) {
	fmt.Printf("pass: "+format+"\n", args...)
}

func goEnv(key string) string {
	out, err := exec.Command("go", "env", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func isGenericGOOS(osName string) bool {
	switch strings.ToLower(strings.TrimSpace(osName)) {
	case "linux", "darwin", "windows":
		return true
	default:
		return false
	}
}

func readOSRelease() (id, versionID string) {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return "", ""
	}
	defer f.Close()

	var m = map[string]string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		m[k] = strings.Trim(v, `"`)
	}
	return m["ID"], m["VERSION_ID"]
}

func detectArchWithVariant() string {
	goarch := strings.ToLower(goEnv("GOARCH"))
	if goarch == "" {
		return ""
	}
	if goarch == "arm" {
		goarm := strings.TrimSpace(goEnv("GOARM"))
		if goarm != "" {
			return "armv" + goarm
		}
		return "arm"
	}
	if goarch == "arm64" {
		// go env 并不提供 v7/v8 这类子版本，这里仅做尽力推断。
		// 大多数 arm64 环境等价于 arm64v8。
		unameOut, err := exec.Command("uname", "-m").Output()
		if err == nil {
			m := strings.ToLower(strings.TrimSpace(string(unameOut)))
			if m == "aarch64" || m == "arm64" {
				return "arm64v8"
			}
		}
		return "arm64"
	}
	return goarch
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
