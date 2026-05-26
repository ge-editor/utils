package utils

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
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

func StartWalkAsync(root string, depth int) (chan FileEvent, *context.CancelFunc) {
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
}

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

// ファイルシステムに存在しないファイルは、
// パス名が一致するかで判定
func SameFile(path1, path2 string) bool {
	// 両方 Stat
	_, err1 := os.Stat(path1)
	_, err2 := os.Stat(path2)

	// 両方存在する場合
	if err1 == nil && err2 == nil {
		// シンボリックリンク解決
		p1, err := filepath.EvalSymlinks(path1)
		if err != nil {
			p1 = path1
		}
		p2, err := filepath.EvalSymlinks(path2)
		if err != nil {
			p2 = path2
		}

		// 絶対パス化
		p1, _ = filepath.Abs(p1)
		p2, _ = filepath.Abs(p2)

		i1, e1 := os.Stat(p1)
		i2, e2 := os.Stat(p2)
		if e1 == nil && e2 == nil {
			return os.SameFile(i1, i2)
		}
		return false
	}

	// 片方だけ存在する
	if (err1 == nil) != (err2 == nil) {
		return false
	}

	// 両方存在しない場合
	if errors.Is(err1, os.ErrNotExist) && errors.Is(err2, os.ErrNotExist) {
		p1 := normalizePath(path1)
		p2 := normalizePath(path2)
		return p1 == p2
	}

	// その他のエラー（権限等）
	return false
}

func normalizePath(p string) string {
	p = filepath.Clean(p)

	abs, err := filepath.Abs(p)
	if err == nil {
		p = abs
	}

	// Windows は大文字小文字を区別しない
	if runtime.GOOS == "windows" {
		p = strings.ToLower(p)
	}

	return p
}

/*
func SameFile(path1, path2 string) bool {
	info1, err1 := os.Stat(path1)
	info2, err2 := os.Stat(path2)

	// 両方存在する場合
	if err1 == nil && err2 == nil {
		return os.SameFile(info1, info2)
	}

	// 片方だけ存在する場合
	if (err1 == nil) != (err2 == nil) {
		return false
	}

	// 両方とも存在しない場合
	if errors.Is(err1, os.ErrNotExist) && errors.Is(err2, os.ErrNotExist) {
		p1 := filepath.Clean(path1)
		p2 := filepath.Clean(path2)
		return p1 == p2
	}

	// その他のエラー（権限など）
	return false
}
*/
