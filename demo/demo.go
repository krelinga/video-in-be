package demo

import (
	"flag"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
)

var demoFlag = flag.String("demo", "", "name of the demo to run")

func Run() {
	demos := map[string]func() error{
		"search_movies": searchMovies,
		"movie_nfo":     movieNfo,
		"refresh":       refresh,
	}
	listDemos := func() string {
		return "available demos: " + strings.Join(slices.Collect(maps.Keys(demos)), ", ")
	}
	if *demoFlag == "" {
		fmt.Fprintf(os.Stderr, "No --demo name provided, %s\n", listDemos())
		os.Exit(1)
	}
	fn, ok := demos[*demoFlag]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown demo %q, %s\n", *demoFlag, listDemos())
		os.Exit(1)
	}
	if err := fn(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
