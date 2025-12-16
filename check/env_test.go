package check

import (
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestGetOS(t *testing.T) {
	got := strings.ToLower(strings.TrimSpace(GetOS()))
	t.Logf("Detected OS: %q", got)
	if got == "" {
		t.Fatalf("GetOS() returned empty")
	}

	goos := strings.ToLower(runtime.GOOS)
	if goos != "linux" {
		if got != goos && !strings.HasPrefix(got, goos) {
			t.Fatalf("GetOS() mismatch on non-linux: got=%q want=%q (or prefix)", got, goos)
		}
		return
	}

	// On Linux, GetOS() may return the generic GOOS ("linux") or a distro descriptor
	// like "ubuntu22.04" (derived from /etc/os-release).
	if got == "linux" {
		return
	}
	if got == "darwin" || got == "windows" {
		t.Fatalf("GetOS() unexpected on linux host: got=%q", got)
	}
}

func TestGetArch(t *testing.T) {
	got := strings.ToLower(strings.TrimSpace(GetArch()))
	t.Logf("Detected Arch: %q", got)
	if got == "" {
		t.Fatalf("GetArch() returned empty")
	}

	goarch := strings.ToLower(runtime.GOARCH)
	switch goarch {
	case "amd64":
		out, err := exec.Command("go", "env", "GOAMD64").Output()
		if err != nil {
			if got != "amd64" {
				t.Fatalf("GetArch() mismatch on amd64 host (no GOAMD64): got=%q want=%q", got, "amd64")
			}
			return
		}
		goamd64 := strings.ToLower(strings.TrimSpace(string(out)))
		if goamd64 == "" {
			if got != "amd64" {
				t.Fatalf("GetArch() mismatch on amd64 host (empty GOAMD64): got=%q want=%q", got, "amd64")
			}
			return
		}
		want := "amd64" + goamd64
		if got != want {
			t.Fatalf("GetArch() mismatch on amd64 host: got=%q want=%q", got, want)
		}
		return
	case "arm64":
		if got != "arm64v8" {
			t.Fatalf("GetArch() mismatch on arm64 host: got=%q want=%q", got, "arm64v8")
		}
	case "arm":
		if !strings.HasPrefix(got, "arm") {
			t.Fatalf("GetArch() mismatch on arm host: got=%q want prefix %q", got, "arm")
		}
	default:
		if got != goarch {
			t.Fatalf("GetArch() mismatch: got=%q want=%q", got, goarch)
		}
	}
}

func TestGetComplier(t *testing.T) {
	got := strings.TrimSpace(GetComplier())
	t.Logf("Detected Compiler Version: %q", got)
	if got == "" {
		t.Fatalf("GetComplier() returned empty")
	}
	// Examples: "go1.24.1", "devel go1.25-..."
	if !strings.HasPrefix(got, "go") && !strings.HasPrefix(got, "devel") {
		t.Fatalf("GetComplier() returned unexpected version string: %q", got)
	}
}
