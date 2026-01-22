package explorer

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
)

func TestExplorr(t *testing.T) {
	fsys := fstest.MapFS{
		"path/to/dir/file0":                         &fstest.MapFile{},
		"path/to/dir/.git/config":                   &fstest.MapFile{},
		"path/to/dir/.gitignore":                    &fstest.MapFile{Data: []byte(".wt/\n")},
		"path/to/dir/.wt/test-1/.git":               &fstest.MapFile{},
		"path/to/dir/.wt/test-1/.git/config":        &fstest.MapFile{},
		"path/to/dir/pkg/foo/path/file1":            &fstest.MapFile{},
		"path/to/dir/pkg/foo/package.json":          &fstest.MapFile{},
		"path/to/dir/pkg/bar/path/file2":            &fstest.MapFile{},
		"path/to/dir/node_modules/baz/package.json": &fstest.MapFile{},
	}

	tests := []struct {
		name       string
		depth      int
		parent     int
		rootFiles  [][]string
		parentDirs [][]string
		ignoreDirs []string
		baseDir    string
		want       []string
	}{
		{
			name:      "Detects root directory (.git/config)",
			depth:     3,
			parent:    2,
			rootFiles: [][]string{{".git", "config"}},
			baseDir:   "/path/to/dir",
			want:      []string{"path/to/dir"},
		},
		{
			name:      "Detects root directory (.git file)",
			depth:     3,
			parent:    2,
			rootFiles: [][]string{{".git"}},
			baseDir:   "/path/to/dir/.wt/test-1",
			want:      []string{"path/to/dir/.wt/test-1"},
		},
		{
			name:      "Explore parent directories to find root directories",
			depth:     1,
			parent:    1,
			rootFiles: [][]string{{".git", "config"}},
			baseDir:   "/path/to/dir/pkg/foo/path",
			want:      []string{"path/to/dir"},
		},
		{
			name:      "Detects root directory (package.json)",
			depth:     3,
			parent:    2,
			rootFiles: [][]string{{"package.json"}},
			baseDir:   "/path/to/dir/pkg/foo",
			want:      []string{"path/to/dir/pkg/foo"},
		},
		{
			name:      "Detects root directory (.git/config and package.json)",
			depth:     3,
			parent:    2,
			rootFiles: [][]string{{".git", "config"}, {"package.json"}},
			baseDir:   "/path/to/dir",
			want:      []string{"path/to/dir", "path/to/dir/node_modules/baz", "path/to/dir/pkg/foo"},
		},
		{
			name:       "Ignore node_modules directory",
			depth:      3,
			parent:     2,
			rootFiles:  [][]string{{".git", "config"}, {"package.json"}},
			ignoreDirs: []string{"node_modules"},
			baseDir:    "/path/to/dir",
			want:       []string{"path/to/dir", "path/to/dir/pkg/foo"},
		},
		{
			name:       "Include directory that has 'pkg' parent directory",
			depth:      3,
			parent:     2,
			rootFiles:  [][]string{{".git", "config"}},
			parentDirs: [][]string{{"pkg"}},
			baseDir:    "/path/to/dir/pkg/foo",
			want:       []string{"path/to/dir", "path/to/dir/pkg/bar", "path/to/dir/pkg/foo"},
		},
		{
			name:      "Ignore gitignored directory",
			depth:     3,
			parent:    2,
			rootFiles: [][]string{{".git", "config"}},
			baseDir:   "/path/to/dir",
			want:      []string{"path/to/dir"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			e := New(fsys, tt.depth, tt.parent, tt.rootFiles, tt.parentDirs, tt.ignoreDirs)
			got, err := e.ExploreRoots(ctx, tt.baseDir)
			if err != nil {
				t.Errorf("%s: %v", tt.name, err)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("%s: mismatch (-want +got):\n%s", tt.name, diff)
			}
		})
	}
}
