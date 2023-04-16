package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/deanveloper/findregex"
)

// useful regexes:
// wildcard selectors: \*[a-zA-Z ():\-_\.]*\{
// before/after: (:before)|(:after)
// ($|\s|,)(\w+)(\s+)(\{|,|^|\+|\[)(\n|^)

var regexMode = flag.Bool("r", false, "lets -f and -x take a `regex` instead of a list of strings")
var includePtr = flag.String("f", "*", "comma-separated file names to `include`")
var excludePtr = flag.String("x", "", "comma-separated path names to `exclude` or skip")
var help = flag.Bool("h", false, "`help`")

func main() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [options...] <regex>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *help {
		flag.Usage()
		fmt.Println()
		fmt.Println("details:")
		fmt.Println(" * the filter provided into -f is only ran on files, not directories")
		fmt.Println(" * the filter provided into -x is ran on both directories and files")
		fmt.Println(" * both -f and -x can use wildcard notation (ie '*.css'), and each value is prepended with an implicit '**/'")
		fmt.Println(" * when in regex mode (-r), -f and -x take Go regular expressions instead of comma-separated lists.")
		fmt.Println(" * regular expressions in Go have most modern regex features: https://github.com/google/re2/wiki/Syntax")
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	regex, err := regexp.Compile(flag.Arg(0))
	if err != nil {
		log.Fatalln(err)
	}

	var filter findregex.FilePathFilterer
	if *regexMode {
		inclRegex, err := regexp.Compile(*includePtr)
		if err != nil {
			log.Fatalln(err)
		}
		exclRegex, err := regexp.Compile(*excludePtr)
		if err != nil {
			log.Fatalln(err)
		}
		filter = findregex.RegexInclExclFilter{
			IncludedFiles: inclRegex,
			ExcludedPaths: exclRegex,
		}
	} else {
		filter = findregex.GlobInclExclFilter{
			IncludedFiles: prefixEach(strings.Split(*includePtr, ","), "**/"),
			ExcludedPaths: prefixEach(strings.Split(*excludePtr, ","), "**/"),
		}
	}

	matcher := findregex.RegexpLineMatcher(*regex)
	results := findregex.Search(".", filter, &matcher)
	for result := range results {
		if result.Err != nil {
			log.Printf("error while searching %q:\n", result.Path)
			log.Println(result.Err)
			os.Exit(2)
		}
		fmt.Printf("%s\t%d\t%s\n", result.Path, result.LineNumber, result.LineText)
	}
}

// returns same slice
func prefixEach(slice []string, prefix string) []string {
	for i, each := range slice {
		slice[i] = prefix + each
	}
	return slice
}
