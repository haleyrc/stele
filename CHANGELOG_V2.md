# v2 Changelog

## New Content Features

### Notes System (NEW)
- Living documents organized by tags instead of chronologically
- Markdown-based with frontmatter (title, tags, pinned)
- Individual note pages at `/notes/{slug}`
- Notes index page with pinned notes at top
- Note tag browsing system (`/notes/tags`, `/notes/tags/{tag}`)

### Post Series (NEW)
- Multi-part post organization via directory structure
- Series metadata via `index.yaml` files
- Automatic chronological ordering within series
- Series navigation sidebar on post pages
- Series index pages at `/{series-slug}`
- Posts accessible at `/posts/{series-slug}/{post-slug}`

### About Page (NEW)
- Optional Markdown-based About page (`about.md`)
- Replaces the old HTML-based pages system
- Social media link integration (GitHub, LinkedIn)

### Archive System (NEW)
- Automatic year-based post archives
- Archive index page listing all years (`/archive`)
- Individual year pages (`/archive/{year}`)

### Enhanced Tag System
- Automatic tag index generation (no longer needs manual configuration)
- Separate tag systems for posts and notes
- Post tags: `/tags`, `/tags/{tag}`
- Note tags: `/notes/tags`, `/notes/tags/{tag}`

## Removed Features
- HTML-based static pages system (replaced by About page)
- Manual menu configuration (now auto-generated from content)

## Development Experience Improvements

### Live Reload (NEW)
- Automatic browser refresh when files change
- WebSocket-based implementation
- Monitors: posts, notes, templates, configuration files
- Smart caching for fast reloads

### Enhanced Development Server
- Configurable port via `--port` flag (default: 3000, was 8081)
- `--live` flag to exclude drafts (simulate production)
- Proper HTTP routing instead of manual refresh endpoint
- Better error handling and logging

### Improved CLI
- ASCII art header with version information
- Clearer command structure and help text
- Better build output configuration (`--out` flag)
- More informative logging during build/dev

## Architecture & Code Quality

### Complete Restructure
- Moved from root-level files to proper `internal/` package structure
- New packages:
  - `internal/compiler/` - Static asset compilation
  - `internal/server/` - Development server with live reload
  - `internal/site/` - Core domain types (Site, Post, Note, Series, etc.)
  - `internal/template/` - Template rendering system
  - `internal/templx/` - Template utilities
  - `internal/textutil/` - Text processing utilities
- Moved `cmd/stele/main.go` to root `main.go`
- Comprehensive test coverage with golden file testing

### Configuration Changes
- Config file renamed: `config.yml` → `stele.yaml`
- Added social media links configuration
- Removed manual menu configuration (auto-generated)

### Build & Tooling
- Updated Go requirement: 1.22 → 1.25
- Enhanced `bin/check` with more thorough validation
- New `bin/clean` script for generated file cleanup
- New `bin/dev` script for development workflows
- Updated GoReleaser configuration

## Documentation
- Extensive README updates with deployment guides
- Platform-specific deployment examples (Cloudflare Pages, AWS S3, Netlify, GitHub Pages, Vercel)
- Detailed file structure documentation
- Comprehensive feature descriptions

## v1 vs v2 Feature Comparison

| Feature | v1 | v2 |
|---------|----|----|
| **Content Types** | Posts only | Posts, Notes, Series |
| **Static Pages** | HTML-based pages/ directory | Markdown About page |
| **Tag System** | Posts only | Posts and Notes separately |
| **Archives** | None | Year-based automatic archives |
| **Series Support** | None | Full series with navigation |
| **Dev Server** | Manual refresh endpoint | Automatic live reload |
| **Port** | Fixed 8081 | Configurable (default 3000) |
| **Navigation** | Manual menu config | Auto-generated from content |
| **Config File** | config.yml | stele.yaml |
| **Social Links** | None | GitHub, LinkedIn support |
| **Architecture** | Root-level files | Proper internal/ packages |
| **Testing** | Basic tests | Comprehensive with golden files |
| **CLI** | Basic commands | Enhanced with version info |
| **Draft Handling** | Dev only | Dev or live mode toggle |
