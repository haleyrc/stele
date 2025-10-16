# About Link Conditional Rendering - Implementation Plan

## Problem Statement

The "About" navigation link is currently rendered unconditionally in both desktop and mobile navigation, even when there is no content to display. This results in a 404 error when users click the link but no `about.md` file exists and no social links are configured.

## Current Behavior

- About link always appears in navigation
- `about.html` only generated when `Site.About != nil` (i.e., when `about.md` exists)
- Server returns 404 when `Site.About == nil`
- Social links are ignored when determining if About page should exist

## Desired Behavior

The About link should only appear when there is actual content to show:

1. **IF** `about.md` exists AND social links exist → render full About page
2. **IF** `about.md` exists AND social links DO NOT exist → render just markdown content
3. **IF** `about.md` DOES NOT exist AND social links DO exist → render just social links
4. **IF** `about.md` DOES NOT exist AND social links DO NOT exist → do NOT show About link, do NOT generate page

## Implementation Plan

### 1. Add Helper Method to Site Struct

**File:** `internal/site/site.go`

**Location:** After the `Site` struct definition (around line 43)

**Add method:**
```go
// HasAboutContent returns true if the site has content for an About page.
// This includes either an about.md file or configured social links.
func (s *Site) HasAboutContent() bool {
	return s.About != nil || s.Config.Social.GitHub != "" || s.Config.Social.LinkedIn != ""
}
```

**Purpose:** Centralize the logic for determining if About page should be shown.

---

### 2. Update About Template

**File:** `internal/template/pages/about.templ`

**Current signature (line 9):**
```go
templ About(s *site.Site, about site.About) {
```

**New signature:**
```go
templ About(s *site.Site, about *site.About) {
```

**Update template body (lines 10-30):**
```go
<article>
	if about != nil {
		<div class="markdown">
			@templ.Raw(about.Content)
		</div>
	}
	if s.Config.Social.GitHub != "" || s.Config.Social.LinkedIn != "" {
		if about != nil {
			<div class="mt-6 pt-6 border-t flex justify-center gap-4">
		} else {
			<div class="flex justify-center gap-4">
		}
			if s.Config.Social.GitHub != "" {
				<a href={ templ.URL(s.Config.Social.GitHub) } target="_blank" rel="noopener noreferrer" class="inline-flex items-center gap-2 hover:underline">
					@icons.GitHub(6)
					<span>GitHub</span>
				</a>
			}
			if s.Config.Social.LinkedIn != "" {
				<a href={ templ.URL(s.Config.Social.LinkedIn) } target="_blank" rel="noopener noreferrer" class="inline-flex items-center gap-2 hover:underline">
					@icons.LinkedIn(6)
					<span>LinkedIn</span>
				</a>
			}
		</div>
	}
</article>
```

**Key changes:**
- Wrap markdown content in `if about != nil` check
- Conditionally apply top border/margin to social links section only when markdown content exists above it

**After editing:** Run `templ generate` to regenerate Go code.

---

### 3. Update Layout Template - Desktop Navigation

**File:** `internal/template/layout.templ`

**Location:** Lines 73-75 (desktop nav)

**Current code:**
```go
<a class="pl-2 hover:underline" href={ templ.URL("/about") }>
	about
</a>
```

**New code:**
```go
if site.HasAboutContent() {
	<a class="pl-2 hover:underline" href={ templ.URL("/about") }>
		about
	</a>
}
```

---

### 4. Update Layout Template - Mobile Navigation

**File:** `internal/template/layout.templ`

**Location:** Lines 110-112 (mobile nav in footer)

**Current code:**
```go
<a class="border-b py-2 hover:bg-gray-100 text-center" href={ templ.URL("/about") }>
	about
</a>
```

**New code:**
```go
if site.HasAboutContent() {
	<a class="border-b py-2 hover:bg-gray-100 text-center" href={ templ.URL("/about") }>
		about
	</a>
}
```

**After editing:** Run `templ generate` to regenerate Go code.

---

### 5. Update Template Renderer

**File:** `internal/template/renderer.go`

**Location:** Line 60 (RenderAbout method)

**Current code:**
```go
func (r *TemplateRenderer) RenderAbout(ctx context.Context, w io.Writer, s *site.Site, about *site.About) error {
	return Layout("About", s, pages.About(s, *about)).Render(ctx, w)
}
```

**New code:**
```go
func (r *TemplateRenderer) RenderAbout(ctx context.Context, w io.Writer, s *site.Site, about *site.About) error {
	return Layout("About", s, pages.About(s, about)).Render(ctx, w)
}
```

**Change:** Remove dereference operator (`*about` → `about`) since template now accepts pointer.

---

### 6. Update Compiler

**File:** `internal/compiler/compiler.go`

**Location:** Lines 132-141 (renderAboutToFile method)

**Current code:**
```go
func (c *Compiler) renderAboutToFile(ctx context.Context, dir string) error {
	if c.Site.About == nil {
		return nil
	}

	path := filepath.Join(dir, "about.html")
	return c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
		return c.Renderer.RenderAbout(ctx, w, c.Site, c.Site.About)
	})
}
```

**New code:**
```go
func (c *Compiler) renderAboutToFile(ctx context.Context, dir string) error {
	if !c.Site.HasAboutContent() {
		return nil
	}

	path := filepath.Join(dir, "about.html")
	return c.renderToFile(ctx, path, func(ctx context.Context, w *os.File) error {
		return c.Renderer.RenderAbout(ctx, w, c.Site, c.Site.About)
	})
}
```

**Change:** Replace `c.Site.About == nil` with `!c.Site.HasAboutContent()`.

**Effect:** About page will be generated when social links exist, even if `about.md` doesn't.

---

### 7. Update Server Handler

**File:** `internal/server/server.go`

**Location:** Lines 73-84 (HandleAbout method)

**Current code:**
```go
func (s *Server) HandleAbout(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	if site.About == nil {
		s.Handle404(w, r)
		return
	}

	s.renderHTML(w, r, func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderAbout(ctx, w, site, site.About)
	})
}
```

**New code:**
```go
func (s *Server) HandleAbout(w http.ResponseWriter, r *http.Request) {
	site := SiteFromContext(r.Context())
	if !site.HasAboutContent() {
		s.Handle404(w, r)
		return
	}

	s.renderHTML(w, r, func(ctx context.Context, w io.Writer) error {
		return s.Renderer.RenderAbout(ctx, w, site, site.About)
	})
}
```

**Change:** Replace `site.About == nil` with `!site.HasAboutContent()`.

**Effect:** Server will render About page when social links exist, even if `about.md` doesn't.

---

## Testing Strategy

### Manual Testing Scenarios

Test all four content combinations:

1. **With about.md + With social links**
   - About link should appear
   - Page shows both markdown content and social links
   - Social links section has top border/margin

2. **With about.md + Without social links**
   - About link should appear
   - Page shows only markdown content
   - No social links section

3. **Without about.md + With social links**
   - About link should appear
   - Page shows only social links
   - No top border/margin on social section

4. **Without about.md + Without social links**
   - About link should NOT appear
   - Navigating to `/about` returns 404
   - No `about.html` generated during build

### Test both modes:
- **Dev server:** `stele serve`
- **Static build:** `stele build` (check if `about.html` exists in output)

---

## Implementation Order

1. Add `HasAboutContent()` method to Site struct
2. Update About template signature and logic
3. Run `templ generate` to regenerate About template
4. Update Layout template (both desktop and mobile nav)
5. Run `templ generate` to regenerate Layout template
6. Update template renderer (remove dereference)
7. Update compiler condition
8. Update server handler condition
9. Test all four scenarios in dev mode
10. Test all four scenarios with build

---

## Files Changed Summary

1. `internal/site/site.go` - Add `HasAboutContent()` method
2. `internal/template/pages/about.templ` - Update signature and conditional rendering
3. `internal/template/pages/about_templ.go` - Auto-generated by templ
4. `internal/template/layout.templ` - Add conditionals around About links
5. `internal/template/layout_templ.go` - Auto-generated by templ
6. `internal/template/renderer.go` - Remove dereference operator
7. `internal/compiler/compiler.go` - Update condition
8. `internal/server/server.go` - Update condition

---

## Edge Cases Handled

- Empty social link strings are treated as "not configured"
- About can be nil (no about.md) but page still renders with social links
- No visual gap when only social links are shown (no unnecessary top margin)
- Consistent behavior between dev server and compiled static site
