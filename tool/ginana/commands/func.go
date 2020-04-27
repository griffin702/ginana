package commands

import (
	"bytes"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

func modPath(p string) string {
	dir := filepath.Dir(p)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			content, _ := ioutil.ReadFile(filepath.Join(dir, "go.mod"))
			mod := RegexpReplace(`module\s+(?P<name>[\S]+)`, string(content), "$name")
			name := strings.TrimPrefix(filepath.Dir(p), dir)
			name = strings.TrimPrefix(name, string(os.PathSeparator))
			if name == "" {
				return fmt.Sprintf("%s/", mod)
			}
			return fmt.Sprintf("%s/%s/", mod, name)
		}
		parent := filepath.Dir(dir)
		if dir == parent {
			return ""
		}
		dir = parent
	}
}

// RegexpReplace replace regexp
func RegexpReplace(reg, src, temp string) string {
	var result []byte
	pattern := regexp.MustCompile(reg)
	for _, subMatches := range pattern.FindAllStringSubmatchIndex(src, -1) {
		result = pattern.ExpandString(result, temp, src, subMatches)
	}
	return string(result)
}

//go:generate packr2
func create() (err error) {
	box := packr.New("all", "./../template/ginana")
	if err = os.MkdirAll(f.Path, 0755); err != nil {
		return
	}
	for _, name := range box.List() {
		if f.ModPrefix != "" && name == "go.mod.tmpl" {
			continue
		}
		tmpl, _ := box.FindString(name)
		i := strings.LastIndex(name, string(os.PathSeparator))
		if i > 0 {
			dir := name[:i]
			if err = os.MkdirAll(filepath.Join(f.Path, dir), 0755); err != nil {
				return
			}
		}
		if strings.HasSuffix(name, ".tmpl") {
			name = strings.TrimSuffix(name, ".tmpl")
		}
		if err = write(filepath.Join(f.Path, name), tmpl); err != nil {
			return
		}
	}
	if err = generate("./..."); err != nil {
		return
	}
	if err = generate("./internal/wire/wire.go"); err != nil {
		return
	}
	return
}

func generate(path string) error {
	cmd := exec.Command("go", "generate", "-x", path)
	cmd.Dir = f.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func write(path, tpl string) (err error) {
	data, err := parse(tpl)
	if err != nil {
		return
	}
	return ioutil.WriteFile(path, data, 0644)
}

func parse(s string) ([]byte, error) {
	t, err := template.New("").Parse(s)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
