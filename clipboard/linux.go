//go:build linux

package clipboard

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func CopyFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("file://%s", absPath)

	if isCommandAvailable("wl-copy") {
		cmd := exec.Command("wl-copy", "--type", "text/uri-list")
		cmd.Stdin = strings.NewReader(uri)
		return cmd.Run()
	}

	if isCommandAvailable("xclip") {
		cmd := exec.Command("xclip", "-selection", "clipboard", "-t", "text/uri-list")
		cmd.Stdin = strings.NewReader(uri)
		return cmd.Run()
	}

	return fmt.Errorf("install 'wl-copy' or 'xclip'")
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
