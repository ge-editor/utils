package utils

import (
	"os"
	"path/filepath"
	"testing"
)

// fsIgnoresCase asks the file system directly, without the code under test:
// it creates a file in dir and looks it up with the case swapped. The file
// stays in dir (the directory of a test is removed with the test).
func fsIgnoresCase(t *testing.T, dir string) bool {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "casetest-probe"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := os.Stat(filepath.Join(dir, "CASETEST-PROBE"))
	return err == nil
}

func TestIsCaseInsensitivePath(t *testing.T) {
	dir := t.TempDir()
	want := fsIgnoresCase(t, dir)

	if got := IsCaseInsensitivePath(dir); got != want {
		t.Errorf("directory with an entry: got %v, want %v", got, want)
	}
	// A file is judged by the directory that holds it.
	if got := IsCaseInsensitivePath(filepath.Join(dir, "casetest-probe")); got != want {
		t.Errorf("file: got %v, want %v", got, want)
	}
	// A path that does not exist is judged by the nearest existing directory.
	if got := IsCaseInsensitivePath(filepath.Join(dir, "no", "such", "path")); got != want {
		t.Errorf("missing path: got %v, want %v", got, want)
	}
	// An empty directory has no entry to test: its own name is used.
	empty := filepath.Join(dir, "Empty")
	if err := os.Mkdir(empty, 0o755); err != nil {
		t.Fatal(err)
	}
	if got := IsCaseInsensitivePath(empty); got != want {
		t.Errorf("empty directory: got %v, want %v", got, want)
	}
}

func TestIsCaseInsensitivePathNothingToExamine(t *testing.T) {
	// A relative path is made absolute, never an error.
	_ = IsCaseInsensitivePath("")
	_ = IsCaseInsensitivePath(".")
}

func TestIsCaseInsensitivePathDoesNotCreateFiles(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Name.txt"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	IsCaseInsensitivePath(dir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "Name.txt" {
		t.Errorf("the directory was changed: %v", entries)
	}
}

func TestProbeEntriesSkipsNamesWithoutLetters(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "123"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, conclusive := probeEntries(dir); conclusive {
		t.Error("a name without letters cannot be tested")
	}
	if err := os.WriteFile(filepath.Join(dir, "Abc"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, conclusive := probeEntries(dir); !conclusive {
		t.Error("a name with letters must give an answer")
	}
}

func TestSwapCase(t *testing.T) {
	if got := swapCase("Foo-bar_1"); got != "fOO-BAR_1" {
		t.Errorf("swapCase = %q", got)
	}
	if got := swapCase("123-_"); got != "123-_" {
		t.Errorf("swapCase of a name without letters changed it: %q", got)
	}
}
