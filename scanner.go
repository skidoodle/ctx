package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

type Scanner struct {
	Root      string
	GlobalCfg *Config
	LocalCfg  *Config
}

func NewScanner(root string, globalCfg *Config) *Scanner {
	s := &Scanner{
		Root:      root,
		GlobalCfg: globalCfg,
	}

	gitIgnorePath := filepath.Join(root, ".gitignore")
	if data, err := os.ReadFile(gitIgnorePath); err == nil {
		s.LocalCfg = parseConfig(string(data))
	}
	return s
}

func (s *Scanner) Scan() ([]string, error) {
	var files []string

	err := filepath.WalkDir(s.Root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(s.Root, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		name := d.Name()
		isDir := d.IsDir()

		if isDir && name == ".git" {
			return filepath.SkipDir
		}

		if s.GlobalCfg.IsIgnored(name, relPath, isDir) {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		}

		if s.LocalCfg != nil && s.LocalCfg.IsIgnored(name, relPath, isDir) {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		}

		if !isDir {
			files = append(files, relPath)
		}
		return nil
	})

	slices.Sort(files)

	return files, err
}
