# findregex

`findregex` is a package that allows you to find regular expressions within files in a fast and efficient manner. This is tested against long directories, along with many nested subdirectories.

This also allows for pruning directories in a much less confusing fashion than GNU `find`, which was actually the purpose behind this in the first place.

# CLI Installation

Go to [releases](https://github.com/deanveloper/findregex/releases) and download the correct binary for your platform, and put the binary somewhere in your `PATH` (ie in `/usr/bin`)

## Install from source

1. Install [Go](https://go.dev/)
2. add `$GOBIN` to your `$PATH`
3. `go get github.com/deanveloper/findregex/findregex`
4. you are ready to Go :)

## Usage

| Flag | Default | Meaning | Example |
| ---- | ------- | ------- | ------- |
| `-f` | `*` | What files to include | `-f '*.css,*.js'` |
| `-x` | (empty) | What directories/files to exclude | `-x 'node_modules,dist,target,*.min.*'` |
| `-r` | false | Change `-f` and `-x` to take regular expressions instead of comma-separated globs | `-r -f 'year_\d{2,4}.js$' -x '(node_modules|dist)$'` |
| `-h` | false | Get into the nitty-gritty | `findregex -h` |

## Examples

 * Find lines which set properties on `window` in javascript files, except in `node_modules`:
   * `findregex -f '*.js' -x 'node_modules' 'window.\S+\s='`

Combine it with other tools for something even more powerful:
 * Find string literals that aren't constants in your Java code
   * `findregex -f '*.java' '".*"' | grep -v 'String \w+ = ".*"'`

## Use as a library

This is not only a CLI command, but a library! To use it as a library, simply run `go get github.com/deanveloper/findregex`, and you can use it as normal. You can view documentation on https://pkg.go.dev/github.com/deanveloper/findregex
