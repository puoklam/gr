package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/puoklam/gr/prompt"
)

// testing promptui
// https://stackoverflow.com/questions/53306447/how-do-i-unit-test-this-promptui-package-written-in-golang
type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb ClosingBuffer) Close() error {
	return nil
}

func TestRootCmd(t *testing.T) {
	path, err := os.MkdirTemp("", "root-cmd")
	srcDir = "./testdata/root/src/"
	defer os.RemoveAll(path)
	assertNoErr(t, err)

	r := ClosingBuffer{
		bytes.NewBufferString("1\n"),
	}
	prompt.SetStdin(r)
	err = run([]string{path})
	assertNoErr(t, err)
}

func assertNoErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
