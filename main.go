package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/haleyrc/stele/internal/compiler"
	"github.com/haleyrc/stele/internal/server"
	"github.com/haleyrc/stele/internal/site"
	"github.com/haleyrc/stele/internal/template"
)

// Version information set by GoReleaser during build.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	ctx := context.Background()

	printHeader()

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
		os.Exit(0)
	default:
		printUsage()
		os.Exit(1)
	}
}

func exitWithError(err error) {
	log.Printf("ERR: %v", err)
	os.Exit(1)
}

func printUsage() {
	fmt.Fprintln(os.Stderr, strings.TrimSpace(usage))
	fmt.Fprintln(os.Stderr)
}

func runBuild(ctx context.Context) {
	buildFlags := flag.NewFlagSet("build", flag.ExitOnError)
	outDir := buildFlags.String("out", "dist", "Output directory for build")
	notesExperiment := buildFlags.Bool("notes-experiment", false, "Enable experimental notes feature")
	if err := buildFlags.Parse(os.Args[2:]); err != nil {
		exitWithError(err)
	}

	site, err := site.New(".", site.SiteOptions{
		IncludeDrafts:   false,
		NotesExperiment: *notesExperiment,
	})
	if err != nil {
		exitWithError(err)
	}

	renderer := template.NewTemplateRenderer()
	compiler := compiler.NewCompiler(renderer, site)
	if err := compiler.Compile(ctx, *outDir, "."); err != nil {
		exitWithError(err)
	}
}

func runDev(ctx context.Context) {
	devFlags := flag.NewFlagSet("dev", flag.ExitOnError)
	port := devFlags.String("port", "3000", "Port to listen on")
	live := devFlags.Bool("live", false, "Exclude draft posts (live mode)")
	notesExperiment := devFlags.Bool("notes-experiment", false, "Enable experimental notes feature")
	if err := devFlags.Parse(os.Args[2:]); err != nil {
		exitWithError(err)
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cache, err := server.NewSiteCache(".", site.SiteOptions{
		IncludeDrafts:   !*live,
		NotesExperiment: *notesExperiment,
	})
	if err != nil {
		exitWithError(err)
	}

	renderer := template.NewTemplateRenderer()

	liveReloader, err := server.NewLiveReloader(*port, renderer, cache)
	if err != nil {
		exitWithError(err)
	}

	liveReloader.Start(ctx)

	if err := liveReloader.ListenAndServe(ctx); err != nil {
		exitWithError(err)
	}
}

func printHeader() {
	fmt.Fprint(os.Stderr, strings.TrimPrefix(logo, "\n"))
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "A no frills blogging platform for people with analysis paralysis")
	fmt.Fprintln(os.Stderr)

	if version == "dev" {
		fmt.Fprintln(os.Stderr, "THIS IS A DEVELOPMENT BUILD. SATISFACTION NOT GUARANTEED.")
	} else {
		fmt.Fprintln(os.Stderr, "Commit: ", commit)
		fmt.Fprintln(os.Stderr, "Date:   ", date)
		fmt.Fprintln(os.Stderr, "Version:", version)
	}
	fmt.Fprintln(os.Stderr)
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
USAGE
  stele COMMAND [OPTIONS]

COMMANDS
  build      Compile static assets for deployment
  dev        Run a development server for previewing content
  help       Print this help message
  version    Print version information

BUILD OPTIONS
  --out               Output directory (default: "dist")
  --notes-experiment  Enable experimental notes feature (default: false)

DEV OPTIONS
  --port              Port to listen on (default: "3000")
  --live              Exclude draft posts (default: false)
  --notes-experiment  Enable experimental notes feature (default: false)
`
