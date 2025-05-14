package main

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
)

func runDemo() error {
	demos := map[string]func() error{
		"search_movies": searchMovies,
	}
	listDemos := func() string {
		return "available demos: " + strings.Join(slices.Collect(maps.Keys(demos)), ", ")
	}

	if len(os.Args) < 2 {
		return errors.New("no demo name provided, " + listDemos())
	}

	demoName := os.Args[1]
	if demo, ok := demos[demoName]; ok {
		return demo()
	} else {
		return fmt.Errorf("unknown demo %q, %s", demoName, listDemos())
	}
}

func main() {
	if err := runDemo(); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
