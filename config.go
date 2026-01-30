package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	configDirName  = "ctx"
	ignoreFileName = "ignore"
)

const defaultIgnoreText = `
.git
.svn
.hg
.gitignore
.idea
.vscode
.vs
.settings
.classpath
.project
.DS_Store
*.swp
*.swo
*.tmp
*.bak
node_modules
bower_components
jspm_packages
vendor
.venv
venv
env
.env
__pycache__
.tox
.mypy_cache
.pytest_cache
.npm
.yarn
.pnpm-store
dist
build
out
target
bin
obj
cmake-build-*
package-lock.json
yarn.lock
pnpm-lock.yaml
bun.lockb
go.sum
Cargo.lock
poetry.lock
Pipfile.lock
*.exe
*.dll
*.so
*.dylib
*.test
*.class
*.jar
*.war
*.ear
*.o
*.obj
*.sqlite
*.db
*.sqlitedb
*.zip
*.tar
*.tar.gz
*.tgz
*.rar
*.7z
*.jpg
*.jpeg
*.png
*.gif
*.ico
*.svg
*.webp
*.mp3
*.mp4
*.mov
*.avi
*.pdf
*.doc
*.docx
*.pyc
*.pyo
*.pyd
*.egg-info
*_templ.go
`

type Config struct {
	ExactIgnores map[string]struct{}
	ExtIgnores   map[string]struct{}
	ComplexGlobs []string
}

func parseConfig(text string) *Config {
	cfg := &Config{
		ExactIgnores: make(map[string]struct{}),
		ExtIgnores:   make(map[string]struct{}),
	}

	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "*.") && strings.Count(line, "*") == 1 && !strings.Contains(line, "/") {
			ext := strings.TrimPrefix(line, "*")
			cfg.ExtIgnores[ext] = struct{}{}
		} else if !strings.ContainsAny(line, "*?[]") && !strings.Contains(line, "/") {
			cfg.ExactIgnores[line] = struct{}{}
		} else {
			cfg.ComplexGlobs = append(cfg.ComplexGlobs, line)
		}
	}
	return cfg
}

func (c *Config) IsIgnored(name, relPath string, isDir bool) bool {
	if _, ok := c.ExactIgnores[name]; ok {
		return true
	}

	if !isDir {
		ext := filepath.Ext(name)
		if _, ok := c.ExtIgnores[ext]; ok {
			return true
		}
	}

	checkPath := filepath.ToSlash(relPath)

	for _, p := range c.ComplexGlobs {
		cleanPattern := strings.TrimSuffix(p, "/")

		if matched, _ := filepath.Match(cleanPattern, name); matched {
			return true
		}

		if matched, _ := filepath.Match(cleanPattern, checkPath); matched {
			return true
		}

		if strings.HasPrefix(checkPath, cleanPattern+"/") {
			return true
		}
	}
	return false
}

func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(configDir, configDirName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, ignoreFileName), nil
}

func loadConfig() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte(defaultIgnoreText), 0644); err != nil {
			return nil, err
		}
		return parseConfig(defaultIgnoreText), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseConfig(string(data)), nil
}

func openConfigFile() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, _ = loadConfig()
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}
