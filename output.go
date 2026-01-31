package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func writeOutput(root string, files []string, outputPath string) (count int64, err error) {
	f, err := os.Create(outputPath)
	if err != nil {
		return 0, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	bw := bufio.NewWriterSize(f, 1024*1024)
	defer func() {
		if ferr := bw.Flush(); ferr != nil && err == nil {
			err = ferr
		}
	}()

	tc := &TokenCounter{w: bw}

	tc.Printf("Project Path: %s\n\n", filepath.Base(root))
	tc.Println("Source Tree:")
	tc.Println("")

	tc.Println("```txt")
	tc.Println(filepath.Base(root))

	if err := writeTree(tc, files); err != nil {
		return 0, err
	}

	tc.Println("```")
	tc.Println("")

	for _, file := range files {
		if file == outputPath || filepath.Base(file) == outputPath {
			continue
		}

		fullPath := filepath.Join(root, file)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ctx: warning: skipping %s: %v\n", file, err)
			continue
		}

		ext := strings.TrimPrefix(filepath.Ext(file), ".")
		if ext == "" {
			ext = "txt"
		}

		tc.Printf("`%s`:\n\n", file)
		tc.Printf("```%s\n", ext)

		if _, err := tc.Write(content); err != nil {
			return 0, err
		}

		if len(content) > 0 && content[len(content)-1] != '\n' {
			if err := tc.WriteByte('\n'); err != nil {
				return 0, err
			}
		}
		tc.Println("```")
		tc.Println("")
	}

	return tc.Count, tc.Err
}

func writeTree(w io.Writer, files []string) error {
	root := make(map[string]any)
	for _, f := range files {
		parts := strings.Split(filepath.ToSlash(f), "/")
		current := root
		for _, part := range parts {
			if _, exists := current[part]; !exists {
				current[part] = make(map[string]any)
			}
			current = current[part].(map[string]any)
		}
	}

	return printNode(w, root, "")
}

func printNode(w io.Writer, node map[string]any, prefix string) error {
	keys := make([]string, 0, len(node))
	for k := range node {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for i, key := range keys {
		isLast := i == len(keys)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		if _, err := fmt.Fprintf(w, "%s%s%s\n", prefix, connector, key); err != nil {
			return err
		}

		children := node[key].(map[string]any)
		if len(children) > 0 {
			childPrefix := prefix + "│   "
			if isLast {
				childPrefix = prefix + "    "
			}
			if err := printNode(w, children, childPrefix); err != nil {
				return err
			}
		}
	}
	return nil
}
