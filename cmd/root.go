package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/puoklam/gr/files"
	"github.com/puoklam/gr/prompt"
	"github.com/spf13/cobra"
)

var ErrInvalidSrcDir = errors.New("invalid source directory")
var ErrInvalidSrcFile = errors.New("invalid source file")
var ErrInvalidDst = errors.New("invalid destination")
var ErrNoSrc = errors.New("source not provided")
var ErrNoDst = errors.New("destination not provided")
var ErrDstNotDir = errors.New("destination not a directory")

var rootCmd = &cobra.Command{
	Use:   "gr-cli",
	Short: "A simple template generator",
	Long:  "gr is a cli written in Go to generate templates with user defined variables",
	RunE:  runCmd,
}

// validate command arguments and flags
func validate(args []string) error {
	if len(args) < 1 {
		return ErrNoSrc
	}
	if len(args) < 2 {
		return ErrNoDst
	}
	if ok, err := files.IsDir(args[1]); err != nil {
		return err
	} else if !ok {
		return ErrDstNotDir
	}
	return nil
}

// command entry point
func runCmd(cmd *cobra.Command, args []string) error {
	return run(args)
}

func run(args []string) error {
	if err := validate(args); err != nil {
		return err
	}
	src, dst := args[0], args[1]
	vars, filemap, err := files.ScanDir(src)
	if err != nil {
		return err
	}
	for _, v := range vars {
		s, err := prompt.Run(fmt.Sprintf("Variable %s to be replaced with", v.Name), nil)
		if err != nil {
			return err
		}
		v.Replace = s
		v.Temp = randString(8)
	}
	return files.Generate(src, dst, filemap)
}

func Exec() error {
	return rootCmd.Execute()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.Grow(n)
	for i, c, r := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if r == 0 {
			c, r = src.Int63(), letterIdxMax
		}
		if j := int(c & letterIdxMask); j < len(letterBytes) {
			sb.WriteByte(letterBytes[j])
			i--
		}
		c >>= letterIdxBits
		r--
	}
	return sb.String()
}
