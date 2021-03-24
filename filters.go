package findregex

import (
	"regexp"

	"github.com/bmatcuk/doublestar/v2"
)

// FilePathFilterer is an interface that determines if a paths and files should be included.
type FilePathFilterer interface {
	PathFilterer
	FileFilterer
}

// PathFilterer is an interface that determines if a path should be searched or skipped.
type PathFilterer interface {
	FilterPath(path string) bool
}

// FileFilterer is an interface that determines if a file should be included or not.
type FileFilterer interface {
	FilterFile(path string) bool
}

// make sure that our provided filters match the interfaces
var _ FilePathFilterer = RegexInclExclFilter{}

// GlobInclExclFilter is a filter that uses glob strings to include files matched by
// IncludedFiles, and prune paths matched by ExcludedPaths.
type GlobInclExclFilter struct {
	IncludedFiles []string
	ExcludedPaths []string
}

func (f GlobInclExclFilter) FilterFile(path string) bool {
	var include bool
	for _, included := range f.IncludedFiles {
		matches, _ := doublestar.PathMatch(included, path)
		if matches {
			include = true
			break
		}
	}
	for _, excluded := range f.ExcludedPaths {
		matches, _ := doublestar.PathMatch(excluded, path)
		if matches {
			include = false
			break
		}
	}

	return include
}

func (f GlobInclExclFilter) FilterPath(path string) bool {
	var exclude bool
	for _, excluded := range f.ExcludedPaths {
		matches, _ := doublestar.PathMatch(excluded, path)
		if matches {
			exclude = true
			break
		}
	}

	return !exclude
}

// RegexInclExclFilter is a filter that uses regular expressions to include files matched by
// IncludedFiles, and prune paths matched by ExcludedPaths.
type RegexInclExclFilter struct {
	IncludedFiles *regexp.Regexp
	ExcludedPaths *regexp.Regexp
}

func (f RegexInclExclFilter) FilterFile(path string) bool {
	return f.IncludedFiles.MatchString(path)
}

func (f RegexInclExclFilter) FilterPath(path string) bool {
	return !f.ExcludedPaths.MatchString(path)
}

// LineMatcher is an interface that checks if a line should be included.
type LineMatcher interface {
	Match(line string) bool
}

type RegexpLineMatcher regexp.Regexp

func (m *RegexpLineMatcher) Match(line string) bool {
	return (*regexp.Regexp)(m).MatchString(line)
}
