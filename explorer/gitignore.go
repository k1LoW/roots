package explorer

import (
	"bufio"
	"bytes"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

type gitignoreCache struct {
	fsys     fs.FS
	repoRoot string
	mu       sync.RWMutex
	patterns map[string][]gitignore.Pattern
	matchers map[string]gitignore.Matcher
}

func newGitignoreCache(fsys fs.FS, repoRoot string) *gitignoreCache {
	return &gitignoreCache{
		fsys:     fsys,
		repoRoot: repoRoot,
		patterns: make(map[string][]gitignore.Pattern),
		matchers: make(map[string]gitignore.Matcher),
	}
}

func (c *gitignoreCache) matcherForDir(dir string) gitignore.Matcher {
	dir = cleanRelDir(dir)
	c.mu.RLock()
	if matcher, ok := c.matchers[dir]; ok {
		c.mu.RUnlock()
		return matcher
	}
	c.mu.RUnlock()

	patterns := c.patternsForDir(dir)
	matcher := gitignore.NewMatcher(patterns)

	c.mu.Lock()
	c.matchers[dir] = matcher
	c.mu.Unlock()

	return matcher
}

func (c *gitignoreCache) patternsForDir(dir string) []gitignore.Pattern {
	dir = cleanRelDir(dir)
	c.mu.RLock()
	if patterns, ok := c.patterns[dir]; ok {
		c.mu.RUnlock()
		return patterns
	}
	c.mu.RUnlock()

	var parentPatterns []gitignore.Pattern
	if dir != "" {
		parentPatterns = c.patternsForDir(parentDir(dir))
	}
	currentPatterns := c.readPatterns(dir)
	combined := append(append([]gitignore.Pattern{}, parentPatterns...), currentPatterns...)

	c.mu.Lock()
	c.patterns[dir] = combined
	c.mu.Unlock()

	return combined
}

func (c *gitignoreCache) readPatterns(dir string) []gitignore.Pattern {
	dir = cleanRelDir(dir)
	gitignorePath := filepath.Join(append([]string{c.repoRoot, dir}, ".gitignore")...)
	data, err := fs.ReadFile(c.fsys, gitignorePath)
	if err != nil {
		return nil
	}

	domain := splitPath(dir)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var patterns []gitignore.Pattern
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
			continue
		}
		patterns = append(patterns, gitignore.ParsePattern(line, domain))
	}
	return patterns
}

func cleanRelDir(dir string) string {
	if dir == "." {
		return ""
	}
	dir = filepath.Clean(dir)
	if dir == "." {
		return ""
	}
	return dir
}

func parentDir(dir string) string {
	if dir == "" {
		return ""
	}
	parent := filepath.Dir(dir)
	if parent == "." {
		return ""
	}
	return parent
}

func splitPath(path string) []string {
	path = strings.Trim(filepath.ToSlash(path), "/")
	if path == "" || path == "." {
		return nil
	}
	return strings.Split(path, "/")
}
