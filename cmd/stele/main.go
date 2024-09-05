package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/haleyrc/stele"
)

const (
	commit  = ""
	date    = ""
	version = ""
)

func main() {
	ctx := context.Background()

	if len(os.Args) == 1 {
		printUsage()
		os.Exit(0)
	}

	command := strings.ToLower(os.Args[1])
	switch command {
	case "build":
		runBuild(ctx)
	case "dev":
		runDev(ctx)
	case "help":
		printUsage()
		os.Exit(0)
	case "version":
		printVersion()
		os.Exit(0)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, strings.TrimSpace(usage))
	fmt.Fprintln(os.Stderr)
}

func printVersion() {
	fmt.Fprint(os.Stderr, logo)
	fmt.Fprintln(os.Stderr, "stele: A no frills blogging platform for people with analysis paralysis")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Commit: ", commit)
	fmt.Fprintln(os.Stderr, "Date:   ", date)
	fmt.Fprintln(os.Stderr, "Version:", version)
}

func runBuild(ctx context.Context) {
	if err := stele.Build(ctx, ".", "dist"); err != nil {
		exitWithError(err)
	}
}

func runDev(ctx context.Context) {
	if err := stele.Build(ctx, ".", "build"); err != nil {
		exitWithError(err)
	}

	http.HandleFunc("GET /refresh", func(w http.ResponseWriter, r *http.Request) {
		if err := stele.Build(ctx, ".", "build"); err != nil {
			log.Println("ERR:", err)
			http.Error(w, "Rebuild failed", http.StatusInternalServerError)
			return
		}
		log.Println("Rebuilt")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	buildDir := os.DirFS(filepath.Join(".", "build"))
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		var fn string
		switch r.URL.Path {
		case "/":
			fn = "index.html"
		case "/rss.xml", "/manifest.webmanifest":
			fn = r.URL.Path
		default:
			fn = r.URL.Path + ".html"
		}
		log.Println("GET", r.URL.Path, fn)
		http.ServeFileFS(w, r, buildDir, fn)
	})

	log.Println("Listening on http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Println("ERR:", err)
		os.Exit(1)
	}
}

func exitWithError(err error) {
	log.Println("ERR:", err)
	os.Exit(1)
}

const logo = `
     _       _
    | |     | |
 ___| |_ ___| | ___
/ __| __/ _ \ |/ _ \
\__ \ ||  __/ |  __/
|___/\__\___|_|\___|
`

const usage = `
A no frills blogging platform for people with analysis paralysis.

USAGE
  stele COMMAND

COMMANDS
  build    Create a set of deployable assets
  dev      Run a local server for previewing content
  help     Print this help message
`
