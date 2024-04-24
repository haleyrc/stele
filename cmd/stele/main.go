package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/haleyrc/stele"
	"github.com/haleyrc/stele/simple"
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
		runBuild(ctx, simple.New(), os.Args[2:]...)
	case "dev":
		runDev(ctx, simple.New())
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	// TODO
}

func runBuild(ctx context.Context, t stele.Theme, args ...string) {
	b, err := stele.New()
	if err != nil {
		exitWithError(err)
	}

	fs := flag.NewFlagSet("build", flag.ExitOnError)
	baseURL := fs.String("b", "", "The base URL of the deployed blog")
	fs.Parse(args)
	fmt.Println(*baseURL)
	if *baseURL != "" {
		b.Config.BaseURL = *baseURL
	}

	if err := stele.Build(ctx, b, t); err != nil {
		exitWithError(err)
	}
}

func runDev(ctx context.Context, t stele.Theme) {
	b, err := stele.New()
	if err != nil {
		exitWithError(err)
	}
	b.Config.BaseURL = "http://localhost:8081"

	if err := stele.Build(ctx, b, t); err != nil {
		exitWithError(err)
	}

	http.HandleFunc("GET /refresh", func(w http.ResponseWriter, r *http.Request) {
		if err := b.Load(); err != nil {
			log.Println("ERR:", err)
			http.Error(w, "Rebuild failed", http.StatusInternalServerError)
			return
		}
		if err := stele.Build(ctx, b, t); err != nil {
			log.Println("ERR:", err)
			http.Error(w, "Rebuild failed", http.StatusInternalServerError)
			return
		}
		log.Println("Rebuilt")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	public := os.DirFS(filepath.Join(".", "public"))
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		var fn string
		switch r.URL.Path {
		case "/":
			fn = "index.html"
		case "/rss.xml":
			fn = r.URL.Path
		case "/manifest.webmanifest":
			fn = r.URL.Path
		default:
			fn = r.URL.Path + ".html"
		}
		log.Println("GET", r.URL.Path, fn)
		http.ServeFileFS(w, r, public, fn)
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
