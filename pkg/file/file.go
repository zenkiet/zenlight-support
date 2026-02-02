package file

import (
	"log/slog"
	"os"
	"path/filepath"
)

type FolderMetrics struct {
	Path         string `json:"path"`
	TotalSize    int64  `json:"totalSize"`
	FileCount    int    `json:"fileCount"`
	LastModified int64  `json:"lastModified"`
}

func GetFolderMetrics(path string) (*FolderMetrics, error) {
	var size int64
	var count int
	var lastMod int64

	cleanPath := filepath.Clean(os.ExpandEnv(path))

	err := filepath.WalkDir(cleanPath, func(path string, dir os.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				slog.Warn("permission denied accessing path", "path", path)
				return filepath.SkipDir
			}
			return filepath.SkipDir
		}

		info, err := dir.Info()
		if err != nil {
			return nil
		}

		if info.ModTime().Unix() > lastMod {
			lastMod = info.ModTime().Unix()
		}

		if !dir.IsDir() {
			size += info.Size()
			count++
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &FolderMetrics{
		Path:         cleanPath,
		TotalSize:    size,
		FileCount:    count,
		LastModified: lastMod,
	}, nil
}
