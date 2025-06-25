package explorer

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
	"sync"
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

// Config holds configuration for concurrent processing
type Config struct {
	MaxWorkers int // worker pool size (default: runtime.NumCPU() * 2)
	BufferSize int // channel buffer size (default: 100)
}

// explorationJob represents a single directory exploration task
type explorationJob struct {
	root  string
	depth int
}

// explorationResult represents the result of a directory exploration
type explorationResult struct {
	roots []string
	err   error
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

// defaultConfig returns default concurrency configuration
func defaultConfig() Config {
	return Config{
		MaxWorkers: runtime.NumCPU() * 2,
		BufferSize: 100,
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
	config := defaultConfig()
	roots, err := e.exploreRootsFromRootConcurrent(ctx, root, depth, config)
	if err != nil {
		return nil, err
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("root not found in %s", baseDir)
	}

	return roots, nil
}

// exploreRootsFromRootConcurrent explores root directories concurrently using worker pool
func (e *Explorer) exploreRootsFromRootConcurrent(ctx context.Context, root string, depth int, config Config) ([]string, error) {
	if depth == 0 || root == "" {
		return nil, nil
	}

	// Check current directory for root files
	var currentRoots []string
	func() {
		for _, rf := range e.rootFiles {
			fp := filepath.Join(append([]string{root}, rf...)...)
			if _, err := fs.Stat(e.fsys, fp); err == nil {
				currentRoots = append(currentRoots, root)
				return
			}
			for _, pd := range e.parentDirs {
				d := filepath.Join(pd...)
				if strings.HasSuffix(filepath.Dir(root), d) {
					if _, err := fs.Stat(e.fsys, filepath.Dir(root)); err == nil {
						currentRoots = append(currentRoots, root)
						return
					}
				}
			}
		}
	}()

	entries, err := fs.ReadDir(e.fsys, root)
	if err != nil {
		return currentRoots, nil // Return current roots even if we can't read subdirectories
	}

	// Filter and collect subdirectories
	var subdirs []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if slices.Contains(e.ignoreDirs, entry.Name()) {
			continue
		}
		subdirs = append(subdirs, filepath.Join(root, entry.Name()))
	}

	if len(subdirs) == 0 {
		return currentRoots, nil
	}
	// Set up worker pool
	workerCount := config.MaxWorkers
	if workerCount <= 0 {
		workerCount = runtime.NumCPU() * 2
	}

	jobs := make(chan explorationJob, config.BufferSize)
	results := make(chan explorationResult, config.BufferSize)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.worker(ctx, jobs, results, config)
		}()
	}

	// Send jobs
	go func() {
		defer close(jobs)
		for _, subdir := range subdirs {
			select {
			case jobs <- explorationJob{root: subdir, depth: depth - 1}:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	var allRoots []string
	allRoots = append(allRoots, currentRoots...)

	for result := range results {
		if result.err != nil {
			return nil, result.err
		}
		allRoots = append(allRoots, result.roots...)
	}

	// Sort results to ensure deterministic output
	sort.Strings(allRoots)
	return allRoots, nil
}

// worker processes exploration jobs
func (e *Explorer) worker(ctx context.Context, jobs <-chan explorationJob, results chan<- explorationResult, config Config) {
	for job := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
			roots, err := e.exploreRootsFromRootConcurrent(ctx, job.root, job.depth, config)
			results <- explorationResult{roots: roots, err: err}
		}
	}
}
