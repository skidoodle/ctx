# ctx

CLI to convert a directory tree and file contents into a single text file for LLM context.
It respects `.gitignore` and comes with sensible defaults for ignoring binaries and lockfiles.

## Installation

### Go Install

```bash
go install github.com/skidoodle/ctx@latest
```

### Binaries

Download pre-compiled binaries for Windows, macOS, and Linux from the [Releases](https://github.com/skidoodle/ctx/releases) page.

## Usage

Generate context for the current directory (outputs to `ctx.txt`):

```bash
ctx .
```

Generate context for a specific folder and save to a custom file:

```bash
ctx -o context.md ./src
```

### Configuration

`ctx` ignores common artifacts (node_modules, .git, binaries) by default.
To edit the global ignore list:

```bash
ctx -config
```

## License
MIT
