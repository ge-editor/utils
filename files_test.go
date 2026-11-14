package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsWithinDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "Proj")
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(base, "ProjX"), 0o755); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		path string
		want bool
	}{
		{dir, true}, // the directory itself
		{filepath.Join(dir, "a.go"), true},
		{filepath.Join(dir, "sub", "b.go"), true},
		{filepath.Join(dir, "new", "c.go"), true}, // does not exist
		{filepath.Join(base, "ProjX", "a.go"), false},
		{base, false},
	}
	for _, c := range cases {
		if got := IsWithinDir(c.path, dir); got != c.want {
			t.Errorf("IsWithinDir(%q, %q) = %v, want %v", c.path, dir, got, c.want)
		}
	}

	// A different case follows the file system of dir.
	other := filepath.Join(base, "PROJ", "a.go")
	if got, want := IsWithinDir(other, dir), IsCaseInsensitivePath(dir); got != want {
		t.Errorf("IsWithinDir with another case = %v, want %v", got, want)
	}
}

func TestRelativePath(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "Proj")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := RelativePath(filepath.Join(dir, "Sub", "File.go"), dir)
	if err != nil || got != filepath.Join("Sub", "File.go") {
		t.Errorf("inside: got %q, %v", got, err)
	}
	got, err = RelativePath(dir, dir)
	if err != nil || got != "." {
		t.Errorf("the directory itself: got %q, %v", got, err)
	}
	got, err = RelativePath(filepath.Join(base, "Other", "X.go"), dir)
	if err != nil || got != filepath.Join("..", "Other", "X.go") {
		t.Errorf("outside: got %q, %v", got, err)
	}

	// Another case of the directory: it is the same directory only on a
	// case-insensitive file system. The result keeps the case of the path.
	got, err = RelativePath(filepath.Join(base, "PROJ", "Sub", "F.go"), dir)
	want := filepath.Join("..", "PROJ", "Sub", "F.go")
	if IsCaseInsensitivePath(dir) {
		want = filepath.Join("Sub", "F.go")
	}
	if err != nil || got != want {
		t.Errorf("another case: got %q, %v, want %q", got, err, want)
	}
}

func TestRestoreCase(t *testing.T) {
	orig := filepath.Join(string(filepath.Separator), "x", "y", "A", "B")
	cases := []struct{ rel, want string }{
		{".", "."},
		{filepath.Join("a", "b"), filepath.Join("A", "B")},
		{filepath.Join("..", "a", "b"), filepath.Join("..", "A", "B")},
		{filepath.Join("..", "..", "b"), filepath.Join("..", "..", "B")},
		{"..", ".."},
	}
	for _, c := range cases {
		if got := restoreCase(c.rel, orig); got != c.want {
			t.Errorf("restoreCase(%q) = %q, want %q", c.rel, got, c.want)
		}
	}
}

func TestIsSameFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "Real.txt")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if !IsSameFile(file, file) {
		t.Error("a file is the same as itself")
	}
	if IsSameFile(file, filepath.Join(dir, "Other.txt")) {
		t.Error("an existing file and a missing one are different")
	}

	// Neither exists: the paths are compared. Another case is the same file
	// only on a case-insensitive file system.
	a := filepath.Join(dir, "Missing.go")
	if !IsSameFile(a, a) {
		t.Error("the same missing path must be the same")
	}
	if IsSameFile(a, filepath.Join(dir, "Missing2.go")) {
		t.Error("different missing paths must differ")
	}
	b := filepath.Join(dir, "MISSING.GO")
	if got, want := IsSameFile(a, b), IsCaseInsensitivePath(dir); got != want {
		t.Errorf("another case: got %v, want %v", got, want)
	}
}
