package ftest

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestAll(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "console.log",
			input:  "console.log('Hello, Godzilla')",
			output: "Hello, Godzilla\n",
		},
		{
			name:   "variable declaration",
			input:  "let foo = 'hello'\nconsole.log(foo)",
			output: "hello\n",
		},
		{
			name:   "assignment",
			input:  "let foo\nfoo = 'hello'\nconsole.log(foo)",
			output: "hello\n",
		},
		{
			name:   "numeric",
			input:  "console.log(10)",
			output: "10\n",
		},
		{
			name:   "binary expression",
			input:  "console.log(1 + 1)",
			output: "2\n",
		},
	}

	bin := filepath.Join(pwd, "..", "bin", "godzilla")
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var out bytes.Buffer
			cmd := exec.Command(bin, "run")
			cmd.Stdin = bytes.NewBufferString(test.input)
			cmd.Stdout = &out
			cmd.Stderr = &out
			if err := cmd.Run(); err != nil {
				t.Fatalf("error running test case %s error=%s stderr=%s", test.name, err, out)
			}

			if want, got := test.output, out.String(); want != got {
				t.Fatalf("output doesn't match for test %s: want=%q got=%q", test.name, want, got)
			}
		})
	}
}
