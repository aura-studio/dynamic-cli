package meta

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"plugin"
	"regexp"
	"strings"
)

// Meta field keys
var Keys = []string{"module", "version", "built", "os", "arch", "compiler", "variant"}

// Result holds extracted metadata
type Result map[string]string

// ReadStrings extracts metadata using the 'strings' command.
// Fast but relies on heuristic pattern matching.
func ReadStrings(soPath string) (Result, error) {
	out, err := exec.Command("strings", soPath).Output()
	if err != nil {
		return nil, fmt.Errorf("strings: %w", err)
	}

	result := make(Result)
	lines := strings.Split(string(out), "\n")

	patterns := map[string]*regexp.Regexp{
		"module":   regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9_.]*\.[a-zA-Z]{2,}/`),
		"version":  regexp.MustCompile(`^v\d+\.\d+\.\d+`),
		"built":    regexp.MustCompile(`^\d{4}-\d{2}-\d{2}_\d{2}:\d{2}:\d{2}_CST`),
		"os":       regexp.MustCompile(`^(linux|darwin|windows)$`),
		"arch":     regexp.MustCompile(`^(amd64|arm64|386)$`),
		"compiler": regexp.MustCompile(`^(gcc|clang|musl-gcc)$`),
		"variant":  regexp.MustCompile(`^(generic|alpine|debian)$`),
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		for key, pattern := range patterns {
			if _, exists := result[key]; !exists && pattern.MatchString(line) {
				result[key] = line
			}
		}
	}

	return result, nil
}

// ReadNm extracts symbol information using 'go tool nm'.
// Returns symbol addresses, not actual values.
func ReadNm(soPath string) (map[string]string, error) {
	out, err := exec.Command("go", "tool", "nm", soPath).Output()
	if err != nil {
		return nil, fmt.Errorf("go tool nm: %w", err)
	}

	symbols := []string{
		"main.MetaDynamicModule",
		"main.MetaDynamicVersion",
		"main.MetaDynamicBuilt",
		"main.MetaToolchainOS",
		"main.MetaToolchainArch",
		"main.MetaToolchainCompiler",
		"main.MetaToolchainVariant",
	}

	lines := strings.Split(string(out), "\n")
	found := make(map[string]string)

	for _, line := range lines {
		for _, sym := range symbols {
			if strings.Contains(line, sym) {
				found[sym] = strings.TrimSpace(line)
			}
		}
	}

	return found, nil
}

// ReadObjdump returns raw .rodata section using 'objdump'.
func ReadObjdump(soPath string) (string, error) {
	out, err := exec.Command("objdump", "-s", "-j", ".rodata", soPath).Output()
	if err != nil {
		return "", fmt.Errorf("objdump: %w", err)
	}
	return string(out), nil
}

// MetaProvider interface for plugin Meta() method
type MetaProvider interface {
	Meta() string
}

// ReadCall loads the .so as a Go plugin and calls Meta().
// Only works with libgo_*.so files (Go plugins).
func ReadCall(soPath string) (Result, string, error) {
	absPath, err := filepath.Abs(soPath)
	if err != nil {
		return nil, "", fmt.Errorf("abs path: %w", err)
	}

	p, err := plugin.Open(absPath)
	if err != nil {
		return nil, "", fmt.Errorf("plugin open: %w", err)
	}

	sym, err := p.Lookup("Tunnel")
	if err != nil {
		return nil, "", fmt.Errorf("lookup Tunnel: %w", err)
	}

	tunnel, ok := sym.(MetaProvider)
	if !ok {
		return nil, "", fmt.Errorf("Tunnel does not implement Meta() method")
	}

	jsonStr := tunnel.Meta()
	result := parseJSON(jsonStr)

	return result, jsonStr, nil
}

func parseJSON(jsonStr string) Result {
	result := make(Result)
	for _, key := range Keys {
		pattern := regexp.MustCompile(`"` + key + `"\s*:\s*"([^"]*)"`)
		if match := pattern.FindStringSubmatch(jsonStr); len(match) > 1 {
			result[key] = strings.TrimSpace(match[1])
		}
	}
	return result
}
