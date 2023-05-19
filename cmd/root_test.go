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
	src := "./testdata/root/src/"
	dst, err := os.MkdirTemp("", "root-cmd")
	defer os.RemoveAll(dst)
	assertNoErr(t, err)

	r := ClosingBuffer{
		bytes.NewBufferString("1\n"),
	}
	prompt.SetStdin(r)
	err = run([]string{src, dst})
	assertNoErr(t, err)
}

func assertNoErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
