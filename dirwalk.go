package utils

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

/*
	 type FileEvent struct {
		Path  string
		IsDir bool
	}
*/
type inodeKey struct {
	dev uint64
	ino uint64
}

const maxWorkers = 16 // ← SSDなら16, HDDなら4〜8推奨

func DirwalkSafe(root string) ([]FileEvent, error) {
	visited := make(map[inodeKey]struct{})
	var result []FileEvent

	err := walk(root, visited, &result)
	return result, err
}

func walk(path string, visited map[inodeKey]struct{}, out *[]FileEvent) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil // 読めない場所は無視
	}

	for _, entry := range entries {
		full := filepath.Join(path, entry.Name())

		linfo, err := os.Lstat(full)
		if err != nil {
			continue
		}

		isDir := false
		var statInfo os.FileInfo

		// symlink？
		if linfo.Mode()&os.ModeSymlink != 0 {
			statInfo, err = os.Stat(full) // ← 実体
			if err != nil {
				continue // 壊れたリンク
			}
		} else {
			statInfo = linfo
		}

		if statInfo.IsDir() {
			isDir = true

			// inode 取得
			sys := statInfo.Sys().(*syscall.Stat_t)
			key := inodeKey{
				dev: uint64(sys.Dev),
				ino: uint64(sys.Ino),
			}

			// 既に訪問済みならスキップ（循環防止）
			if _, ok := visited[key]; ok {
				continue
			}
			visited[key] = struct{}{}

			// 再帰
			walk(full, visited, out)
		}

		*out = append(*out, FileEvent{
			Path:  full,
			IsDir: isDir,
		})
	}

	return nil
}

func DirwalkParallelCtxDepth(ctx context.Context, root string, maxDepth int) []FileEvent {
	var results []FileEvent
	var mu sync.Mutex

	visited := make(map[inodeKey]struct{})
	var visitedMu sync.Mutex

	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	var walk func(string, int)
	walk = func(path string, depth int) {
		defer wg.Done()

		// depth 制限
		if maxDepth >= 0 && depth > maxDepth {
			return
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			return
		}

		entries, err := os.ReadDir(path)
		<-sem
		if err != nil {
			return
		}

		for _, entry := range entries {
			select {
			case <-ctx.Done():
				return
			default:
			}

			full := filepath.Join(path, entry.Name())

			linfo, err := os.Lstat(full)
			if err != nil {
				continue
			}

			var statInfo os.FileInfo
			isDir := false

			if linfo.Mode()&os.ModeSymlink != 0 {
				statInfo, err = os.Stat(full)
				if err != nil {
					continue
				}
			} else {
				statInfo = linfo
			}

			if statInfo.IsDir() {
				isDir = true

				sys := statInfo.Sys().(*syscall.Stat_t)
				key := inodeKey{dev: uint64(sys.Dev), ino: uint64(sys.Ino)}

				visitedMu.Lock()
				if _, ok := visited[key]; ok {
					visitedMu.Unlock()
					continue
				}
				visited[key] = struct{}{}
				visitedMu.Unlock()

				wg.Add(1)
				go walk(full, depth+1)
			}

			mu.Lock()
			results = append(results, FileEvent{Path: full, IsDir: isDir})
			mu.Unlock()
		}
	}

	wg.Add(1)
	go walk(root, 0)
	wg.Wait()

	return results
}
