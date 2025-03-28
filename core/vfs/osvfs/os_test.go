package osvfs_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/core/repo"
	"github.com/microsoft/typescript-go/core/tspath"
	"github.com/microsoft/typescript-go/core/vfs/osvfs"
	"gotest.tools/v3/assert"
)

func TestOS(t *testing.T) {
	t.Parallel()

	fs := osvfs.FS()

	goMod := filepath.Join(repo.RootPath, "go.mod")
	goModPath := tspath.NormalizePath(goMod)

	t.Run("ReadFile", func(t *testing.T) {
		t.Parallel()

		expectedRaw, err := os.ReadFile(goMod)
		assert.NilError(t, err)
		expected := string(expectedRaw)

		contents, ok := fs.ReadFile(goModPath)
		assert.Assert(t, ok)
		assert.Equal(t, contents, expected)
	})

	t.Run("Realpath", func(t *testing.T) {
		t.Parallel()

		expected := goModPath
		if runtime.GOOS == "windows" {
			// Windows drive letters can be lowercase, but realpath will always return uppercase.
			expected = strings.ToUpper(expected[:1]) + expected[1:]
		}
		realpath := fs.Realpath(goModPath)
		assert.Equal(t, realpath, expected)
	})

	t.Run("UseCaseSensitiveFileNames", func(t *testing.T) {
		t.Parallel()

		// Just check that it works.
		fs.UseCaseSensitiveFileNames()

		switch runtime.GOOS {
		case "windows":
			assert.Assert(t, !fs.UseCaseSensitiveFileNames())
		case "linux":
			assert.Assert(t, fs.UseCaseSensitiveFileNames())
		}
	})
}
