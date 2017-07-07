package build

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jingweno/godzilla/ast"
	"github.com/jingweno/godzilla/compiler"
	"github.com/jingweno/godzilla/source"
)

func Run(parserPath string, r io.Reader) (string, error) {
	source, err := compileSource(parserPath, r)
	if err != nil {
		return "", err
	}

	main, err := writeMainFile(source)
	if err != nil {
		return "", err
	}

	if err := goFmt(main); err != nil {
		return "", err
	}

	return main, nil
}

func compileSource(parserPath string, r io.Reader) (*source.Code, error) {
	c := exec.Command(parserPath)
	c.Stdin = r
	stdoutStderr, err := c.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error parsing JavaScript %s: %s", err, stdoutStderr)
	}

	f := &ast.File{}
	if err := json.NewDecoder(bytes.NewBuffer(stdoutStderr)).Decode(f); err != nil {
		return nil, err
	}

	return compiler.Compile(f), nil
}

func writeMainFile(code *source.Code) (string, error) {
	mainDir, err := ioutil.TempDir("", "main")
	if err != nil {
		return "", err
	}

	mainFile, err := os.OpenFile(filepath.Join(mainDir, "main.go"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}

	if err := code.WriteTo(mainFile); err != nil {
		return "", err
	}

	if err := mainFile.Close(); err != nil {
		return "", err
	}

	return mainFile.Name(), nil
}

func goFmt(mainFile string) error {
	cmd := exec.Command("go", "fmt", mainFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Errorf("error running `go fmt %s`: error=%s out=%s", err, out)
	}

	return nil
}
