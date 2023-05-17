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

// flags
var (
	srcDir  string
	srcFile string
	// exclude string
)

var ErrInvalidSrcDir = errors.New("invalid source directory")
var ErrInvalidSrcFile = errors.New("invalid source file")
var ErrInvalidDst = errors.New("invalid destination")
var ErrNoDst = errors.New("destination not provided")
var ErrDstNotDir = errors.New("destination not a directory")

var rootCmd = &cobra.Command{
	Use:   "gr-cli",
	Short: "A simple template generator",
	Long:  "gr is a cli written in Go to generate templates with user defined variables",
	RunE:  runE,
}

func init() {
	rootCmd.Flags().StringVarP(&srcDir, "src", "s", "", "Source directory to read from")
	rootCmd.Flags().StringVarP(&srcFile, "file", "f", "", "Source file to read from")
	rootCmd.MarkFlagsMutuallyExclusive("src", "file")
}

// validate command arguments and flags
func validate(args []string) error {
	if len(args) < 1 {
		return ErrNoDst
	}
	if ok, err := files.IsDir(args[0]); err != nil {
		return err
	} else if !ok {
		return ErrDstNotDir
	}
	return nil
}

// command entry point
func runE(cmd *cobra.Command, args []string) error {
	if err := validate(args); err != nil {
		return err
	}
	dst := args[0]
	vars, filemap, err := files.ScanDir(srcDir)
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
	return files.Generate(srcDir, dst, filemap)
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
