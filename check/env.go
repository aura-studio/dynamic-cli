package check

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// GetOS returns a best-effort OS descriptor.
// - On Linux with /etc/os-release: "<id><version_id>" like "ubuntu22.04".
// - Otherwise: returns GOOS like "linux"/"darwin"/"windows".
func GetOS() string {
	out, err := exec.Command("go", "env", "GOOS").Output()
	if err != nil {
		return strings.ToLower(strings.TrimSpace(runtime.GOOS))
	}
	goos := strings.ToLower(strings.TrimSpace(string(out)))
	if goos != "linux" {
		return goos
	}

	f, err := os.Open("/etc/os-release")
	if err != nil {
		return goos
	}
	defer f.Close()

	m := map[string]string{}
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
	id := strings.ToLower(strings.TrimSpace(m["ID"]))
	ver := strings.TrimSpace(m["VERSION_ID"])
	if id == "" || ver == "" {
		return goos
	}
	return id + strings.TrimSpace(ver)
}

// GetArch returns current arch with best-effort variant.
// Examples: amd64, armv7, arm64v8.
func GetArch() string {
	out, err := exec.Command("go", "env", "GOARCH").Output()
	if err != nil {
		goarch := strings.ToLower(strings.TrimSpace(runtime.GOARCH))
		if goarch == "arm64" {
			return "arm64v8"
		}
		return goarch
	}
	goarch := strings.ToLower(strings.TrimSpace(string(out)))
	if goarch == "" {
		return ""
	}
	if goarch == "arm" {
		outArm, errArm := exec.Command("go", "env", "GOARM").Output()
		if errArm != nil {
			return "arm"
		}
		goarm := strings.TrimSpace(string(outArm))
		if goarm != "" {
			return "armv" + goarm
		}
		return "arm"
	}
	if goarch == "arm64" {
		// Go does not expose an arm64 variant; treat as v8 by default.
		return "arm64v8"
	}
	return goarch
}

// GetComplier returns `go env GOVERSION`, e.g. go1.20.5.
func GetComplier() string {
	out, err := exec.Command("go", "env", "GOVERSION").Output()
	if err != nil {
		return strings.TrimSpace(runtime.Version())
	}
	return strings.TrimSpace(string(out))
}
