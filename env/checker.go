package env

import (
	"fmt"
	"os"
	"strings"
)

func CheckOS(targetOS string) bool {
	// 直接比对 GetOS() 的返回值（严格相等）。
	if targetOS == "" {
		return true
	}

	actualOS := strings.ToLower(strings.TrimSpace(GetOS()))
	if actualOS == "" {
		warnf("cannot detect GOOS (target=%s)", targetOS)
		return false
	}

	expected := strings.ToLower(strings.TrimSpace(targetOS))
	if actualOS != expected {
		warnf("OS mismatch (target=%s actual=%s)", expected, actualOS)
		return false
	}
	okf("OS match (target=%s actual=%s)", expected, actualOS)
	return true
}

func CheckArch(targetArch string) bool {
	// 直接比对 GetArch() 的返回值（严格相等）。
	if targetArch == "" {
		return true
	}

	actual := GetArch()
	if actual == "" {
		warnf("cannot detect arch (target=%s)", targetArch)
		return false
	}

	expected := strings.ToLower(strings.TrimSpace(targetArch))
	actualLower := strings.ToLower(actual)

	if expected == actualLower {
		okf("ARCH match (target=%s actual=%s)", expected, actualLower)
		return true
	}
	warnf("ARCH mismatch (target=%s actual=%s)", expected, actualLower)
	return false
}

func CheckCompiler(targetCompiler string) bool {
	// 检查Go的编译器版本，如go1.20.5等，不符合则警告（严格相等）。
	if targetCompiler == "" {
		return true
	}
	actual := GetCompiler()
	if actual == "" {
		warnf("cannot detect GOVERSION (target=%s)", targetCompiler)
		return false
	}
	expected := strings.TrimSpace(targetCompiler)
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
