package toolchain_test

import (
	"strings"
	"testing"

	"github.com/aura-studio/dynamic-cli/toolchain"
)

func TestDetect(t *testing.T) {
	v := toolchain.Detect()
	if strings.TrimSpace(v.OS) == "" {
		t.Fatalf("Detect().OS is empty")
	}
	if strings.TrimSpace(v.Arch) == "" {
		t.Fatalf("Detect().Arch is empty")
	}
	if strings.TrimSpace(v.Compiler) == "" {
		t.Fatalf("Detect().Compiler is empty")
	}
}

func TestDescribe(t *testing.T) {
	if strings.TrimSpace(toolchain.Describe(toolchain.KindOS)) == "" {
		t.Fatalf("Describe(os) is empty")
	}
	if strings.TrimSpace(toolchain.Describe(toolchain.KindArch)) == "" {
		t.Fatalf("Describe(arch) is empty")
	}
	if strings.TrimSpace(toolchain.Describe(toolchain.KindCompiler)) == "" {
		t.Fatalf("Describe(compiler) is empty")
	}
}
