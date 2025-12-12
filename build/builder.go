package build

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

var templateMap = map[string]string{}

type RenderData struct {
	Name        string
	Module      string
	Package     string
	Version     string
	House       string
	Environment string
	Variant     string
	Dir         string
	OS          string
	Arch        string
	Compiler    string
}

type Builder struct {
	config       *RenderData
	user         *user.User
	netRCPath    string
	netRCBakPath string
}

func New(c *RenderData) *Builder {
	user, err := user.Current()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// Determine netrc path cross-platform
	var netrcName string
	if runtime.GOOS == "windows" {
		netrcName = "_netrc"
	} else {
		netrcName = ".netrc"
	}
	return &Builder{
		config:       c,
		user:         user,
		netRCPath:    filepath.Join(user.HomeDir, netrcName),
		netRCBakPath: filepath.Join(user.HomeDir, netrcName+".go_dynamic_bak"),
	}
}

func (b *Builder) Build() {
	fmt.Println("start...")
	defer fmt.Println("done!")

	// Use netrc content from env when provided
	netrc := os.Getenv("DYNAMIC_CLI_NETRC")
	if strings.TrimSpace(netrc) != "" {
		b.bakNetRC()
		b.writeNetRC(netrc)
		defer b.restoreNetRC()
	}

	b.generate()
	b.runBuilder()
}

// bakNetRC backup netrc file if exsits
func (b *Builder) bakNetRC() {
	fmt.Println("bakup", b.netRCPath)
	if _, err := os.Stat(b.netRCPath); err == nil {
		if err := os.Rename(b.netRCPath, b.netRCBakPath); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}
}

// writeNetRC write netrc file from Builder.config.NetRC
func (b *Builder) writeNetRC(netrc string) {
	fmt.Println("write", b.netRCPath)
	f, err := os.Create(b.netRCPath)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := f.WriteString(netrc); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

// restoreNetRC restore netrc file if exsits
func (b *Builder) restoreNetRC() {
	fmt.Println("restore", b.netRCPath)
	if _, err := os.Stat(b.netRCBakPath); err == nil {
		if err := os.Rename(b.netRCBakPath, b.netRCPath); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}
}

// generate generate packages, go.mod and builder.sh
func (b *Builder) generate() {
	for pathTemplateStr, textTemplateStr := range templateMap {
		var pathBuilder strings.Builder
		if pathTemplate, err := template.New("dynamic").Parse(pathTemplateStr); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		} else if err := pathTemplate.Execute(&pathBuilder, b.config); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		path := pathBuilder.String()
		fmt.Println("generate", path)
		if _, err := os.Stat(filepath.Dir(path)); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					fmt.Println("error:", err)
					os.Exit(1)
				}
			} else {
				fmt.Println("error:", err)
				os.Exit(1)
			}
		}

		f, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		defer f.Close()

		if textTemplate, err := template.New("dynamic").Parse(textTemplateStr); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		} else if err := textTemplate.Execute(f, b.config); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}
}

// runBuilder run ./builder.sh
func (b *Builder) runBuilder() {
	if err := os.Chmod(filepath.Join(b.config.Dir, "builder.sh"), 0755); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	cmd := exec.Command("sh", "-c", "./builder.sh")
	cmd.Dir = b.config.Dir
	cmd.Env = append(os.Environ(), "USER="+b.user.Username, "HOME="+b.user.HomeDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// cat builder.sh
	builderPath := filepath.Join(cmd.Dir, "builder.sh")
	builderFile, err := os.Open(builderPath)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	builderContent, err := io.ReadAll(builderFile)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("exec builder.sh\n================\n%s================\n", string(builderContent))
	builderFile.Close()

	if err := cmd.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
