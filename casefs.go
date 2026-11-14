package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// probeEntryLimit is how many directory entries IsCaseInsensitivePath reads
// at most while looking for one that can be tested.
const probeEntryLimit = 32

// IsCaseInsensitivePath reports whether the file system that holds path
// ignores the case of file names.
//
// The answer comes from the file system, not from the operating system: the
// same OS can have both kinds (a USB stick, a network share, a case-sensitive
// directory on Windows or macOS). path may be a file or a directory and may not
// exist; the nearest existing directory above it is examined. Nothing is
// created or changed.
//
// It looks at an entry of that directory, swaps the case of its name and asks
// the file system whether the new name is the same file. A directory without
// such an entry is judged by its own name in its parent. When the file system
// cannot be examined the answer is false (case is significant).
//
// Each call does a few file system calls, so keep the result when it is needed
// repeatedly for the same directory.
func IsCaseInsensitivePath(path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	dir, ok := nearestExistingDir(abs)
	if !ok {
		return false
	}
	if insensitive, conclusive := probeEntries(dir); conclusive {
		return insensitive
	}
	return probeDirName(dir)
}

// nearestExistingDir returns p when it is an existing directory, otherwise
// the nearest existing directory above it.
func nearestExistingDir(p string) (string, bool) {
	for {
		if fi, err := os.Stat(p); err == nil {
			if fi.IsDir() {
				return p, true
			}
			return filepath.Dir(p), true
		}
		parent := filepath.Dir(p)
		if parent == p {
			return "", false
		}
		p = parent
	}
}

// probeEntries tests the file system of dir with the entries inside it. This
// also covers a directory that is the root of a mounted file system. conclusive
// is false when no entry could be tested (the directory is empty, or the names
// have no letters, or the directory cannot be read).
func probeEntries(dir string) (insensitive, conclusive bool) {
	f, err := os.Open(dir)
	if err != nil {
		return false, false
	}
	defer f.Close()

	names, _ := f.Readdirnames(probeEntryLimit) // a partial list is enough
	for _, name := range names {
		swapped := swapCase(name)
		if swapped == name {
			continue
		}
		fi, err := os.Lstat(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		alt, err := os.Lstat(filepath.Join(dir, swapped))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return false, true // the swapped name does not resolve
			}
			continue // not conclusive (permission, name too long, ...)
		}
		// It resolves: the same file means the file system ignores case; a
		// different file means both names exist, so case is significant.
		return os.SameFile(fi, alt), true
	}
	return false, false
}

// probeDirName tests the file system with the name of dir itself, or of the
// nearest parent whose name has a letter.
func probeDirName(dir string) bool {
	cur := dir
	for {
		base := filepath.Base(cur)
		if swapped := swapCase(base); swapped != base {
			fi, err := os.Stat(cur)
			if err != nil {
				return false
			}
			alt, err := os.Stat(filepath.Join(filepath.Dir(cur), swapped))
			return err == nil && os.SameFile(fi, alt)
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return false
		}
		cur = parent
	}
}

// swapCase turns upper case letters to lower case and the other way round.
func swapCase(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case unicode.IsUpper(r):
			return unicode.ToLower(r)
		case unicode.IsLower(r):
			return unicode.ToUpper(r)
		}
		return r
	}, s)
}
