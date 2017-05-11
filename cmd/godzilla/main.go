package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/jingweno/godzilla/ast"
	"github.com/jingweno/godzilla/compiler"
	"github.com/jingweno/godzilla/source"
	"github.com/spf13/cobra"
)

var parserPath string

func main() {
	rootCmd := &cobra.Command{Use: "godzilla"}
	cmdRun := &cobra.Command{
		Use:   "run",
		Short: "compile and run JavaScript program",
		RunE:  runRun,
	}
	rootCmd.AddCommand(cmdRun)
	rootCmd.PersistentFlags().StringVarP(&parserPath, "parser-path", "p", filepath.Join(filepath.Dir(os.Args[0]), "godzilla-parser"), "path to godzilla-parser")
	rootCmd.Execute()
}

func runRun(cmd *cobra.Command, args []string) error {
	c := exec.Command(parserPath)
	c.Stdin = os.Stdin
	stdoutStderr, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error parsing JavaScript %s: %s", err, stdoutStderr)
	}

	m := make(map[string]interface{})
	if err := json.NewDecoder(bytes.NewBuffer(stdoutStderr)).Decode(&m); err != nil {
		return err
	}

	f := &ast.File{}
	f.UnmarshalMap(m)

	code := compiler.Compile(f)
	mainFile, err := writeMainFile(code)
	if err != nil {
		return err
	}

	return goRun(mainFile)
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

func goRun(mainFile string) error {
	goBin, err := exec.LookPath("go")
	if err != nil {
		return err
	}

	return syscall.Exec(goBin, []string{"go", "run", mainFile}, os.Environ())
}
