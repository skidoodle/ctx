package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const defaultOutputFilename = "ctx.txt"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ctx: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Usage = usage
	configFlag := flag.Bool("config", false, "open global ignore file for editing")
	outputFlag := flag.String("o", defaultOutputFilename, "output filename")
	flag.Parse()

	if *configFlag {
		return openConfigFile()
	}

	args := flag.Args()
	if len(args) == 0 {
		usage()
		return nil
	}

	targetDir := args[0]
	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}

	start := time.Now()

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ctx: warning: could not load config, using defaults (%v)\n", err)
		cfg = parseConfig(defaultIgnoreText)
	}

	scanner := NewScanner(absPath, cfg)
	files, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scanning directory: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found in %s", absPath)
	}

	tokenCount, err := writeOutput(absPath, files, *outputFlag)
	if err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	duration := time.Since(start).Round(time.Millisecond)
	fmt.Printf("ctx: generated %s (%d files, ~%d tokens) in %v\n",
		*outputFlag, len(files), tokenCount, duration)

	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: ctx [options] <directory>\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -config      Open global ignore file for editing\n")
	fmt.Fprintf(os.Stderr, "  -o <file>    Output filename (default %q)\n", defaultOutputFilename)
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  ctx .                  # Generate context for current dir\n")
	fmt.Fprintf(os.Stderr, "  ctx -o out.txt src/    # Scan src/ and save to out.txt\n")
	fmt.Fprintf(os.Stderr, "  ctx -config            # Open global ignore file\n")
}
