/*
Copyright © 2025 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/k1LoW/roots/explorer"
	"github.com/k1LoW/roots/version"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	depth                int
	parent               int
	rootFiles            []string
	ignoreDirs           []string
	parentDirs           []string
	defaultRootFilePaths [][]string = [][]string{
		{".git", "config"}, // .git/config is a file that exists in the root directory of a Git project
		{"go.mod"},         // go.mod is a file that exists in the root directory of a Go project
		{"package.json"},   // package.json is a file that exists in the root directory of a Node.js project
		{"Cargo.toml"},     // Cargo.toml is a file that exists in the root directory of a Rust project
	}
	defaultIgnoreDirs []string = []string{
		"node_modules",
		"vendor",
		"testdata",
	}
	fast bool
)

var rootCmd = &cobra.Command{
	Use:          "roots",
	Short:        "roots is a tool for exploring multiple root directories",
	Long:         "roots is a tool for exploring multiple root directories, such as those in a monorepo project.",
	Version:      version.Version,
	SilenceUsage: true,
	Args:         cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		sysRoot := os.Getenv("SystemDrive")
		if sysRoot == "" {
			sysRoot = "/"
		}
		var baseDirs []string
		var rootFilePaths [][]string
		var parentDirPaths [][]string
		if len(rootFiles) == 0 {
			rootFilePaths = defaultRootFilePaths
		} else {
			for _, rf := range rootFiles {
				rootFilePaths = append(rootFilePaths, strings.Split(rf, string(filepath.Separator)))
			}
		}
		if len(parentDirs) > 0 {
			for _, pd := range parentDirs {
				parentDirPaths = append(parentDirPaths, strings.Split(pd, string(filepath.Separator)))
			}
		}

		switch {
		case len(args) == 1:
			abs, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}
			baseDirs = []string{abs}
		case !isatty.IsTerminal(os.Stdin.Fd()):
			b, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			for _, dir := range strings.Split(strings.Trim(string(b), "\n"), "\n") {
				abs, err := filepath.Abs(dir)
				if err != nil {
					return err
				}
				baseDirs = append(baseDirs, abs)
			}
		default:
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			baseDirs = []string{wd}
		}

		e := explorer.New(os.DirFS(sysRoot), depth, parent, rootFilePaths, parentDirPaths, ignoreDirs)
		eg, ctx := errgroup.WithContext(ctx)
		if !fast {
			eg.SetLimit(1)
		}
		for _, baseDir := range baseDirs {
			eg.Go(func() error {
				dirs, err := e.ExploreRoots(ctx, baseDir)
				if err != nil {
					return err
				}
				for _, dir := range dirs {
					fmt.Println(filepath.Join(sysRoot, dir))
				}
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&depth, "depth", "d", 3, "Depth for exploring directories")
	rootCmd.Flags().IntVarP(&parent, "parent", "p", 2, "Number of parent root directories to explore")
	rootCmd.Flags().StringSliceVarP(&rootFiles, "root-file", "", []string{}, "File or directory that exist in the root directory")
	rootCmd.Flags().StringSliceVarP(&ignoreDirs, "ignore-dir", "", defaultIgnoreDirs, "Directory to ignore")
	rootCmd.Flags().StringSliceVarP(&parentDirs, "parent-dir", "", []string{}, "Directory that exists as a parent directory of the root directory")
	rootCmd.Flags().BoolVarP(&fast, "fast", "", false, "Fast mode")
}
