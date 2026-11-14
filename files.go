package utils

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Return last part of path
// Trim right side path separator and split by separator
func LastPartOfPath(path string) string {
	sep := string(os.PathSeparator)
	path = strings.TrimRight(path, sep)
	components := strings.Split(path, sep)
	return components[len(components)-1]
}

// Get the list of files in the directory
func Dirwalk(dir string, depth int) []string {
	if depth < 0 {
		return []string{}
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return []string{}
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			subPaths := Dirwalk(filepath.Join(dir, file.Name()), depth-1)
			paths = append(paths, subPaths...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}

/*
//
ctx, cancel := context.WithCancel(context.Background())

fileCh := make(chan FileEvent, 512)
var wg sync.WaitGroup

// 同時探索数制限（超重要）
sem := make(chan struct{}, 64)

wg.Add(1)
sem <- struct{}{}

go func() {
	defer func() { <-sem }()
	WalkAsync(ctx, root, 999, fileCh, &wg, sem)
}()

go func() {
	wg.Wait()
	close(fileCh)
}()

//
ticker := time.NewTicker(50 * time.Millisecond)

for {
	select {
	case ev := <-fileCh:
		buffer = append(buffer, ev)

	case <-ticker.C:
		files = append(files, buffer...)
		buffer = buffer[:0]
		requestRedraw()
	}
}

*/

type FileEvent struct {
	Path  string
	IsDir bool
}

/* // Get the list of files in the directory
func Dirwalk2(dir string) []FileEvent {
	if dir == "" {
		dir = "."
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return []FileEvent{}
	}

	var paths []FileEvent
	for _, file := range files {
		// symlink 先が dir の場合でも IsDir を true にしたい
		paths = append(paths, FileEvent{Path: file.Name(), IsDir: file.IsDir()})
	}

	return paths
}
*/

func Dirwalk2(dir string) []FileEvent {
	if dir == "" {
		dir = "."
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var paths []FileEvent

	for _, entry := range entries {
		fullpath := filepath.Join(dir, entry.Name())

		isDir := entry.IsDir()

		// symlink の場合だけリンク先を見る
		if entry.Type()&os.ModeSymlink != 0 {
			info, err := os.Stat(fullpath) // ← これが「リンク解決」
			if err == nil && info.IsDir() {
				isDir = true
			}
		}

		paths = append(paths, FileEvent{
			Path:  entry.Name(),
			IsDir: isDir,
		})
	}

	return paths
}

/* func StartWalkAsync(root string, depth int) (chan FileEvent, *context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	fileCh := make(chan FileEvent, 512)
	var wg sync.WaitGroup

	// 同時探索数制限（超重要）
	sem := make(chan struct{}, 64)

	wg.Add(1)
	sem <- struct{}{}

	go func() {
		defer func() { <-sem }()
		WalkAsync(ctx, root, depth, fileCh, &wg, sem)
	}()

	go func() {
		wg.Wait()
		close(fileCh)
	}()

	return fileCh, &cancel
} */

func WalkAsync(
	ctx context.Context,
	root string,
	depth int,
	out chan<- FileEvent,
	wg *sync.WaitGroup,
	sem chan struct{},
) {
	defer wg.Done()

	match := ""
	path := ""
	if root != string(filepath.Separator) && strings.HasSuffix(root, string(filepath.Separator)) {
		path = strings.TrimRight(root, string(filepath.Separator))
	} else {
		match = filepath.Base(root)
	}

	// キャンセル確認
	if ctx != nil {
		select {
		case <-ctx.Done():
			return
		default:
		}
	}

	// 深さ制限
	if depth == 0 {
		return
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	for {
		entries, err := f.Readdir(100)

		for _, e := range entries {

			// キャンセル早期チェック
			if ctx != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
			}

			full := filepath.Join(path, e.Name())

			// out <- FileEvent{full, e.IsDir()}
			if match == "" {
				out <- FileEvent{e.Name(), e.IsDir()}
			} else {
				if ContainsAllCharacters(e.Name(), match) {
					out <- FileEvent{e.Name(), e.IsDir()}
				}
			}

			if e.IsDir() {
				wg.Add(1)

				// goroutine制限
				sem <- struct{}{}

				go func(p string) {
					defer func() { <-sem }()
					WalkAsync(ctx, p, depth-1, out, wg, sem)
				}(full)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
	}
}

func CopyFile(src, dest string) error {
	d, err := os.Create(dest)
	if err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}

	_, err = io.Copy(d, s)
	return err
}

func ExistsFile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// IsWithinDir reports whether path is within dir, including dir itself.
// It compares normalized paths without resolving symbolic links. Case is
// ignored when the file system that holds dir ignores it
// (see IsCaseInsensitivePath).
func IsWithinDir(path, dir string) bool {
	fold := IsCaseInsensitivePath(dir)
	path = normalizePath(path, fold)
	dir = normalizePath(dir, fold)

	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return false
	}

	return rel == "." ||
		(rel != ".." &&
			!strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

// RelativePath returns the path relative to dir.
// It returns an error if the paths cannot be compared. Case is ignored when
// the file system that holds dir ignores it (see IsCaseInsensitivePath); the
// result keeps the case of path.
func RelativePath(path, dir string) (string, error) {
	fold := IsCaseInsensitivePath(dir)
	rel, err := filepath.Rel(normalizePath(dir, fold), normalizePath(path, fold))
	if err != nil || !fold {
		return rel, err
	}
	return restoreCase(rel, normalizePath(path, false)), nil
}

// restoreCase gives rel, which was computed from lower-cased paths, the case
// of the original path: the ".." elements stay and the other elements are the
// last elements of orig.
func restoreCase(rel, orig string) string {
	if rel == "." {
		return rel
	}
	sep := string(filepath.Separator)
	elems := strings.Split(rel, sep)
	origElems := strings.Split(orig, sep)

	var out []string
	keep := 0
	for _, e := range elems {
		if e == ".." {
			out = append(out, e)
		} else {
			keep++
		}
	}
	if keep > len(origElems) {
		return rel
	}
	out = append(out, origElems[len(origElems)-keep:]...)
	return filepath.Join(out...)
}

// IsSameFileInfo reports whether fi1 and fi2 describe the same file.
func IsSameFileInfo(fi1, fi2 os.FileInfo) bool {
	if fi1 == nil || fi2 == nil {
		return false
	}
	return os.SameFile(fi1, fi2)
}

// IsSameFile reports whether path1 and path2 refer to the same file.
// If neither path exists, it compares their normalized paths; case is ignored
// when the file systems of both paths ignore it (see IsCaseInsensitivePath).
func IsSameFile(path1, path2 string) bool {
	fi1, err1 := os.Stat(path1)
	fi2, err2 := os.Stat(path2)

	if err1 == nil && err2 == nil {
		return os.SameFile(fi1, fi2)
	}

	// Both paths do not exist: compare normalized paths.
	if errors.Is(err1, os.ErrNotExist) &&
		errors.Is(err2, os.ErrNotExist) {
		fold := IsCaseInsensitivePath(path1) && IsCaseInsensitivePath(path2)
		return normalizePath(path1, fold) == normalizePath(path2, fold)
	}

	return false
}

// normalizePath makes p absolute and clean. With fold it is also lower-cased,
// so that paths on a case-insensitive file system compare equal.
func normalizePath(p string, fold bool) string {
	p = filepath.Clean(p)

	abs, err := filepath.Abs(p)
	if err == nil {
		p = abs
	}

	if fold {
		p = strings.ToLower(p)
	}

	return p
}
