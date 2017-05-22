package ftest

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestConsoleLog(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	cmd := exec.Command(filepath.Join(pwd, "..", "bin", "godzilla"), "run")
	cmd.Stdin = bytes.NewBufferString(`console.log("Hello, Godzilla")`)
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("error running bin/godzilla error=%s stderr=%s", err, out)
	}

	if want, got := "Hello, Godzilla\n", out.String(); want != got {
		t.Fatalf("output doesn't match: want=%q got=%q", want, got)
	}
}
