package cmd

import (
	"errors"
	"fmt"

	"github.com/puoklam/gr/files"
	"github.com/puoklam/gr/prompt"
	"github.com/spf13/cobra"
)

// flags
var (
	src     string
	srcFile string
	dst     string
	// exclude string
)

var ErrNotDir = errors.New("not a directory")

var rootCmd = &cobra.Command{
	Use:   "gr-cli",
	Short: "A simple template generator",
	Long:  "gr is a cli written in Go to generate templates with user defined variables",
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok, err := files.IsDir(dst); err != nil {
			return err
		} else if !ok {
			return ErrNotDir
		}
		// return files.Generate(src)
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
		}
		return files.Generate(src, dst, filemap)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&src, "src", "s", "", "Source directory to read from")
	rootCmd.Flags().StringVarP(&srcFile, "file", "f", "", "Source file to read from")
	rootCmd.Flags().StringVarP(&dst, "dst", "d", "", "Destination directory to write to")
	rootCmd.MarkFlagsMutuallyExclusive("src", "file")
	rootCmd.MarkFlagsRequiredTogether("src", "dst")
	// rootCmd.MarkFlagsRequiredTogether("file", "dst")
}

func Exec() error {
	return rootCmd.Execute()
}
