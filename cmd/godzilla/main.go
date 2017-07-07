package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/jingweno/godzilla/build"
	"github.com/spf13/cobra"
)

var (
	parserPath     string
	debug          bool
	buildGoOutFile string
)

func main() {
	rootCmd := &cobra.Command{Use: "godzilla"}
	cmdBuild := &cobra.Command{
		Use:   "build",
		Short: "compile JavaScript program",
		RunE:  runBuild,
	}
	cmdBuild.PersistentFlags().StringVarP(&buildGoOutFile, "output", "o", "", "output file")
	cmdRun := &cobra.Command{
		Use:   "run",
		Short: "compile and run JavaScript program",
		RunE:  runRun,
	}
	rootCmd.AddCommand(cmdBuild)
	rootCmd.AddCommand(cmdRun)
	rootCmd.PersistentFlags().StringVarP(&parserPath, "parser-path", "p", filepath.Join(filepath.Dir(os.Args[0]), "godzilla-parser"), "path to godzilla-parser")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "x", false, "print compiled source")
	rootCmd.Execute()
}

func runBuild(cmd *cobra.Command, args []string) error {
	r := os.Stdin
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer f.Close()

		r = f
		if buildGoOutFile == "" {
			buildGoOutFile = strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		}
	}

	mainFile, err := build.Run(parserPath, r)
	if err != nil {
		return err
	}

	if debug {
		err := printCompiledSource(mainFile)
		if err != nil {
			return err
		}
	}

	return goBuild(mainFile, buildGoOutFile)
}

func runRun(cmd *cobra.Command, args []string) error {
	r := os.Stdin
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}

	mainFile, err := build.Run(parserPath, r)
	if err != nil {
		return err
	}

	if debug {
		err := printCompiledSource(mainFile)
		if err != nil {
			return err
		}
	}

	return goRun(mainFile)
}

func printCompiledSource(mainFile string) error {
	out, err := ioutil.ReadFile(mainFile)
	if err != nil {
		return err
	}

	fmt.Printf("%s", out)

	return nil
}

func goBuild(mainFile, outFile string) error {
	goBin, err := exec.LookPath("go")
	if err != nil {
		return err
	}

	return syscall.Exec(goBin, []string{"go", "build", "-o", outFile, mainFile}, os.Environ())
}

func goRun(mainFile string) error {
	goBin, err := exec.LookPath("go")
	if err != nil {
		return err
	}

	return syscall.Exec(goBin, []string{"go", "run", mainFile}, os.Environ())
}
