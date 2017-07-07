package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jingweno/godzilla/build"
	"github.com/spf13/cobra"
)

var (
	parserPath string
)

func main() {
	rootCmd := &cobra.Command{
		Use:  "godzillac",
		RunE: run,
	}
	rootCmd.PersistentFlags().StringVarP(&parserPath, "parser-path", "p", filepath.Join(filepath.Dir(os.Args[0]), "godzilla-parser"), "path to godzilla-parser")
	rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
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

	out, err := ioutil.ReadFile(mainFile)
	if err != nil {
		return err
	}

	fmt.Printf("%s", out)

	return nil
}
