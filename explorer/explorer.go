package explorer

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Explorer struct {
	fsys       fs.FS
	sysRoot    string
	depth      int
	parent     int
	rootFiles  [][]string
	parentDirs [][]string
	ignoreDirs []string
}

func New(fsys fs.FS, depth int, parent int, rootFiles, parentDirs [][]string, ignoreDirs []string) *Explorer {
	sysRoot := os.Getenv("SystemDrive")
	if sysRoot == "" {
		sysRoot = "/"
	}
	return &Explorer{
		fsys:       fsys,
		sysRoot:    sysRoot,
		depth:      depth,
		parent:     parent,
		rootFiles:  rootFiles,
		parentDirs: parentDirs,
		ignoreDirs: ignoreDirs,
	}
}

func (e *Explorer) ExploreRoots(ctx context.Context, baseDir string) ([]string, error) {
	current := strings.TrimLeft(baseDir, e.sysRoot)
	fi, err := fs.Stat(e.fsys, current)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", baseDir)
	}

	// Explore parent root directories
	var root string
	parent := e.parent
	for {
		if current == filepath.Dir(current) {
			break
		}
		func() {
			for _, rf := range e.rootFiles {
				fp := filepath.Join(append([]string{current}, rf...)...)
				if _, err := fs.Stat(e.fsys, fp); err == nil {
					root = current
					parent--
					return
				}
			}
			for _, pd := range e.parentDirs {
				d := filepath.Join(pd...)
				if strings.HasSuffix(filepath.Dir(current), d) {
					if _, err := fs.Stat(e.fsys, filepath.Dir(current)); err == nil {
						root = current
						parent--
						return
					}
				}
			}
		}()
		if parent == 0 {
			break
		}
		current = filepath.Dir(current)
	}

	// Explore child root directories
	depth := e.depth
	root = strings.TrimLeft(root, e.sysRoot)
	roots, err := e.exploreRootsFromRoot(ctx, root, depth)
	if err != nil {
		return nil, err
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("root not found in %s", baseDir)
	}

	return roots, nil
}

func (e *Explorer) exploreRootsFromRoot(ctx context.Context, root string, depth int) ([]string, error) {
	var roots []string
	if depth == 0 || root == "" {
		return nil, nil
	}
	func() {
		for _, rf := range e.rootFiles {
			fp := filepath.Join(append([]string{root}, rf...)...)
			if _, err := fs.Stat(e.fsys, fp); err == nil {
				roots = append(roots, root)
				return
			}
			for _, pd := range e.parentDirs {
				d := filepath.Join(pd...)
				if strings.HasSuffix(filepath.Dir(root), d) {
					if _, err := fs.Stat(e.fsys, filepath.Dir(root)); err == nil {
						roots = append(roots, root)
						return
					}
				}
			}
		}
	}()
	entries, err := fs.ReadDir(e.fsys, root)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if slices.Contains(e.ignoreDirs, entry.Name()) {
			continue
		}
		subRoot := filepath.Join(root, entry.Name())
		subRoots, err := e.exploreRootsFromRoot(ctx, subRoot, depth-1)
		if err != nil {
			return nil, err
		}
		roots = append(roots, subRoots...)
	}

	return roots, nil
}
