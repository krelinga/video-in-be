package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/krelinga/video-in-be/demo"
	"github.com/krelinga/video-in-be/manual"
	"github.com/krelinga/video-in-be/server"
)

var mode string

func init() {
	flag.StringVar(&mode, "mode", "server", "Mode in which the application runs")
}

func main() {
	flag.Parse()
	fmt.Printf("Running in mode: %s\n", mode)

	switch mode {
	case "server":
		server.Start()
	case "demo":
		demo.Run()
	case "manual":
		manual.WriteNfo()
	default:
		fmt.Fprintf(os.Stderr, "Unknown mode: %s, supported modes are 'server', 'demo' and 'manual'\n", mode)
		os.Exit(1)
	}
}
