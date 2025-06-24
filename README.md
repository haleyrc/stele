# `stele`

An opinionated static site generator for people with analysis paralysis.

---

[![Go Reference](https://pkg.go.dev/badge/github.com/haleyrc/stele.svg)](https://pkg.go.dev/github.com/haleyrc/stele)
![Build Status](https://github.com/haleyrc/stele/actions/workflows/test.yml/badge.svg)

## Overview

`stele` is a static-site generator focused on simplicity over configurability. The goal of this project is to create a fast, purpose-built SSG to meet my needs when spinning up a new blog. If you're hoping for a more featureful, I highly recommend you check out the [alternatives](#alternatives) below as well as the broader world of SSG tooling. There's definitely something out there that does what you want.

As for this project, here's what you _do_ get:

* A CLI with dev server and build commands
* Markdown posts with frontmatter
* Post series for organizing related posts
* Markdown notes for living documents
* Optional About page (also in Markdown)
* Automatic archive pages (by year)
* Automatic tag pages (if posts have tags)
* Automatic note tag pages (if notes have tags)
* Automatic web manifest
* Automatic RSS feed
* (Mostly) responsive design

and, perhaps most critically:

* No decisions

If that sounds good enough for your use-case, then do please try it out!

## Install from source

To install from source you will need Go >= 1.25 installed:

```
go install github.com/haleyrc/stele@latest
```

## Usage

Assuming `stele` is already installed and in your path, using the CLI is pretty straightforward.

### Development Server

Running `stele dev` will start up a local development server:

```
stele dev
```

The development server includes automatic live reload - when you save changes to posts, notes, templates, or configuration files, the browser will automatically refresh to show your updates.

Available options:

* `--port` - Port to listen on (default: `3000`)
* `--live` - Exclude draft posts to simulate production (default: `false`)

### Building for Production

Once you are satisfied with your post/page/whatever, you can do a full build by running `stele build`:

```
stele build
```

This will create a `dist/` folder with all of the static assets for the site:

* `index.html` - The home page with recent posts
* `about.html` - Your about page (only if `about.md` exists)
* `notes.html` - Notes index page (only if notes exist)
* `notes/` - Individual note pages (one per note)
* `notes/tags.html` - Note tag index page (only if notes exist)
* `notes/tags/{tag}.html` - Notes for each specific tag (one page per tag)
* `posts/` - Individual post pages (one per post)
* `posts/{series}/` - Individual series post pages (one per post in series)
* `{series}.html` - Series index pages (one per series)
* `archive.html` - Archive index page listing all years
* `archive/{year}.html` - Posts from each specific year (one page per year)
* `tags.html` - Tag index page
* `tags/{tag}.html` - Posts for each specific tag (one page per tag)
* `manifest.webmanifest` - Web app manifest
* `rss.xml` - RSS feed

These assets should be deployable as-is to something like an S3 bucket or you can have your favorite host (e.g. Cloudflare Pages, Netlify, etc.) build and deploy them for you. How to set that up is outside of the scope of this guide, but shouldn't be too difficult for someone with experience on these platforms.

Available options:

* `--out` - Output directory (default: `dist`)

For either of these commands to work correctly, you will need to make sure that your source directory is laid out in the standard `stele` format.

## Deployment

Once you've built your site with `stele build`, the `dist/` directory contains static assets ready for deployment. Here are examples for popular hosting platforms:

### Cloudflare Pages

**Via Git Integration:**

1. Push your project to GitHub/GitLab
2. In Cloudflare Pages dashboard, connect your repository
3. Configure build settings:
   - **Build command:** `stele build`
   - **Output directory:** `dist`

**Via Wrangler CLI:**

```bash
stele build
npx wrangler pages deploy dist --project-name=my-blog
```

### AWS S3 + CloudFront

**S3 bucket setup:**

```bash
# Create bucket
aws s3 mb s3://my-blog

# Enable static website hosting
aws s3 website s3://my-blog --index-document index.html

# Sync files
stele build
aws s3 sync dist/ s3://my-blog --delete
```

**CloudFront** (optional): Create a distribution pointing to your S3 bucket for HTTPS and CDN.

### Netlify

**Create `netlify.toml` in your project root:**

```toml
[build]
  command = "stele build"
  publish = "dist"
```

**Manual deployment:**

```bash
stele build
npx netlify deploy --prod --dir=dist
```

### GitHub Pages

**Create `.github/workflows/deploy.yml`:**

```yaml
name: Deploy

on:
  push:
    branches: [main]

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  deploy:
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - run: go install github.com/haleyrc/stele@latest
      - run: stele build
      - uses: actions/upload-pages-artifact@v3
        with:
          path: ./dist
      - uses: actions/deploy-pages@v4
        id: deployment
```

Enable GitHub Pages in repository settings (Source: "GitHub Actions").

### Vercel

**Create `vercel.json` in your project root:**

```json
{
  "buildCommand": "stele build",
  "outputDirectory": "dist"
}
```

**Manual deployment:**

```bash
stele build
npx vercel --prod --cwd dist
```

### General Notes

- All platforms serve `index.html` files automatically for directory requests
- Configure custom domains in your platform's dashboard
- Most platforms provide automatic HTTPS certificates
- The `dist/` directory is self-contained and portable across any static host

## File structure

The `stele` CLI relies on a specific file structure to work correctly:

```
.
├── stele.yaml
├── about.md (optional)
├── notes/ (optional)
└── posts/
    ├── standalone-post.md
    └── series-name/
        ├── index.yaml
        ├── post-one.md
        └── post-two.md
```

### `stele.yaml`

At the root of your project you will need a YAML file containing the configuration for your blog. The available configuration values are:

* `title` - The name of your blog
  ```
  title: My First Blog
  ```
* `description` - A short description of your blog's content/purpose
  ```
  description: A personal blog about fish
  ```
* `author` - Probably your name
  ```
  author: Grace Hopper
  ```
* `baseURL` - The FQDN and protocol for your blog
  ```
  baseURL: https://myblog.example.com
  ```
* `categories` - A list of categories that describe the content of the blog
  ```
  categories:
    - blog
    - programming
  ```
  > This can also be written in short-form: `categories: [blog, programming]`.
* `social` - Social media links to display on the About page (optional)
  ```
  social:
    github: https://github.com/username
    linkedin: https://www.linkedin.com/in/username
  ```
  > Currently supports `github` and `linkedin`. These links will only be displayed if an `about.md` file exists.

> [!NOTE]
> The site navigation is automatically generated based on your content. If you have posts, an "archive" link appears. If any posts have tags, a "tags" link appears. An RSS feed link is always present.

### `about.md` (optional)

An optional About page can be created at the root of your project as `about.md`. This file should contain Markdown content describing yourself or your blog. If present, the About page will be rendered at `/about.html` and can include social media links configured in the `social` section of `stele.yaml`.

The About page is written in Markdown and will be automatically parsed and rendered with any configured social media links displayed below the content.

### `notes/` (optional)

Notes are living documents organized by tags rather than chronologically. Unlike posts, notes are expected to be modified over time and don't have publication dates.

`stele` will load and parse any files in the `notes/` subdirectory with a `.md` extension. The note filename should be URL-safe as it will be used as the slug (e.g., `golang-tips.md` becomes `/notes/golang-tips.html`).

The note content should follow the format:

```yaml
---
title: Note Title
tags: [tag1, tag2]
pinned: false
---

Your content here.
```

**Required fields:**
* `title` - The note title
* `tags` - A list of tags to associate with the note (can be an empty array)

**Optional fields:**
* `pinned` - Pin this note to the top of the notes index (default: `false`)

**Notes index:**
* The notes index (`/notes.html`) displays pinned notes at the top (sorted alphabetically by title)
* Below pinned notes, a tag index is shown to browse notes by tag
* Individual tag pages (`/notes/tags/{tag}.html`) list all notes with that tag, sorted alphabetically by title

### `posts/`

Posts in `stele`, as with most SSGs, are written in markdown with some minimal frontmatter.

#### Standalone Posts

`stele` will load and parse any files in the `posts/` subdirectory with a `.md` extension. The post filename should be URL-safe as it will be used as the slug for the post in the final assets (e.g., `my-first-post.md` becomes `/posts/my-first-post.html`).

The post content should follow the format:

```yaml
---
title: Your Title Here
description: Your description here.
tags: [tag1, tag2]
draft: false
date: YYYY-MM-DD
---

Your content here.
```

**Required fields:**
* `title` - The post title
* `description` - A short description of the post

**Optional fields:**
* `tags` - A list of tags to associate with the post. If at least one post has tags, `stele` will automatically generate a tags index page and individual tag pages.
* `draft` - Whether this is a draft post (default: `false`). Draft posts are visible in the development server but excluded from production builds.
* `date` - The publication date in `YYYY-MM-DD` format.

**Important notes about dates:**
* Published posts (non-drafts) **must** have a date
* Draft posts **must not** have a date
* When running the development server, drafts are automatically assigned today's date so they appear at the top of the post list

#### Post Series

Posts can be organized into series by creating a subdirectory under `posts/` with an `index.yaml` file. Series are perfect for multi-part tutorials, related topics, or any collection of posts that build on each other.

**Directory structure:**
```
posts/
└── go-basics/
    ├── index.yaml
    ├── variables.md
    ├── functions.md
    └── types.md
```

**Series metadata (`index.yaml`):**
```yaml
name: "Go Basics"
description: "Learn the fundamentals of Go programming"
```

**Required fields:**
* `name` - The display name for the series

**Optional fields:**
* `description` - A description shown on the series index page

**Series behavior:**
* Posts within a series use the same frontmatter format as standalone posts
* Posts are automatically ordered chronologically (oldest first) within the series
* Each series gets an index page at `/{series-slug}` (e.g., `/go-basics`)
* Series posts are accessible at `/posts/{series-slug}/{post-slug}` (e.g., `/posts/go-basics/variables`)
* When viewing a post in a series, a navigation sidebar shows all posts in the series with the current post highlighted
* Series posts are included in the main post feed, archive pages, and tag pages alongside standalone posts
* The series slug is derived from the directory name (e.g., `go-basics/` becomes the slug `go-basics`)

---

That's really all you need to get started. Everything else is handled by the framework!

## Future Improvements

- Homebrew deployment
- Better image handling

## Alternatives

- [Gatsby](https://www.gatsbyjs.com/)
- [Hugo](https://gohugo.io/)
- [Jekyll](https://jekyllrb.com/)