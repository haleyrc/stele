package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/haleyrc/stele"
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
		runBuild(ctx, os.Args[2:]...)
	case "dev":
		runDev(ctx)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	// TODO
}

func runBuild(ctx context.Context, _ ...string) {
	if err := stele.Build(ctx, ".", "dist"); err != nil {
		exitWithError(err)
	}
}

func runDev(ctx context.Context) {
	// TODO
	//
	// b.Config.BaseURL = "http://localhost:8081"
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
		case "/rss.xml":
			fn = r.URL.Path
		case "/manifest.webmanifest":
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
