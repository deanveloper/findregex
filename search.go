package findregex

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// Workers is the number of workers that Search will use when reading files.
var Workers = 5

// SearchResult is a simple struct containing information about a line found by SearchFiles.
type SearchResult struct {
	Path string

	LineText   string
	LineNumber int

	Err error
}

// Search is a combination call to FindFiles and SearchFiles, noteworthy that this function will
// create several workers for SearchFiles.
func Search(path string, filter FilePathFilterer, matcher LineMatcher) <-chan SearchResult {
	files := FindFiles(path, filter)

	chans := make([]<-chan SearchResult, Workers)
	for i := 0; i < Workers; i++ {
		chans[i] = SearchFiles(files, matcher)
	}
	return coalesce(chans)
}

// FindFiles takes a path and returns a string of all files that match the search params
func FindFiles(path string, filter FilePathFilterer) <-chan string {
	files := make(chan string, Workers)
	go func() {
		filepath.Walk(path, func(fpath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				if !filter.FilterPath(fpath) {
					return filepath.SkipDir
				}
			} else {
				if filter.FilterFile(fpath) {
					files <- fpath
				}
			}
			return nil
		})
		close(files)
	}()
	return files
}

// SearchFiles takes a channel of paths, and searches for the regex in each file for the path.
func SearchFiles(files <-chan string, matcher LineMatcher) <-chan SearchResult {
	resultsChan := make(chan SearchResult)
	go func() {
		for filename := range files {
			searchResults := searchFile(filename, matcher)
			for _, result := range searchResults {
				resultsChan <- result
			}
		}
		close(resultsChan)
	}()
	return resultsChan
}
func searchFile(filename string, matcher LineMatcher) []SearchResult {
	file, err := os.Open(filename)
	if err != nil {
		return []SearchResult{{Path: filename, Err: err}}
	}
	defer file.Close()
	results := readerContains(file, matcher.Match)

	for i, result := range results {
		result.Path = filename
		results[i] = result
	}
	return results
}
func readerContains(r io.Reader, containsFunc func(line string) bool) []SearchResult {
	reader := bufio.NewReader(r)
	var lines []SearchResult
	for num := 0; ; num++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			return lines
		}
		if containsFunc(line) {
			lines = append(lines, SearchResult{
				LineNumber: num,
				LineText:   strings.Trim(line, "\n\t {"),
			})
		}
	}
}

func coalesce(chans []<-chan SearchResult) <-chan SearchResult {
	coalesced := make(chan SearchResult)

	selectCases := make([]reflect.SelectCase, len(chans))
	for i, ch := range chans {
		selectCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}

	go func() {
		for len(selectCases) > 0 {
			i, value, ok := reflect.Select(selectCases)
			if !ok {
				selectCases = remove(selectCases, i)
				continue
			}
			coalesced <- value.Interface().(SearchResult)
		}
		close(coalesced)
	}()

	return coalesced
}

func remove(s []reflect.SelectCase, i int) []reflect.SelectCase {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
