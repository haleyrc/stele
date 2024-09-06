# stele

TODO: Logo

A no frills blogging platform for people with analysis paralysis.

---

[![Go Reference](https://pkg.go.dev/badge/github.com/haleyrc/stele.svg)](https://pkg.go.dev/github.com/haleyrc/stele)
![Build Status](https://github.com/haleyrc/stele/actions/workflows/test.yml/badge.svg)

## Overview

`stele` is a static-site generator focused on simplicity over configurability. The goal of this project is to create a fast, purpose-built SSG to meet my needs when spinning up a new blog. If you're hoping for a more featureful, I highly recommend you check out the [alternatives](#alternatives) below as well as the broader world of SSG tooling. There's definitely something out there that does what you want.

As for this project, here's what you _do_ get:

* A CLI
* Markdown posts
* Static pages
* Automatic web manifest
* Automatic RSS feed
* (Pretty) responsive design

and, perhaps most critically:

* No decisions

If that sounds good enough for your use-case, then do please try it out!

## Install from source

To install from source you will need Go >= 1.22 installed:

```
go install github.com/haleyrc/stele
```

## Usage

Assuming `stele` is already installed and in your path, using the CLI is pretty straightforward. Running `stele dev` will start up a local development server. This does a full build of the static assets and makes them available at http://localhost:8081. At the moment, the site will not automatically rebuild on changes, but you can visit http://localhost:8081/refresh to force a rebuild.

Once you are satisfied with your post/page/whatever, you can do a full build by running `stele build`. This will create a `dist/` folder with all of the static assets for the site. These assets should be deployable as-is to something like an S3 bucket or you can have your favorite host (e.g. Cloudflare Pages, Netlify, etc.) build and deploy them for you. How to set that up is outside of the scope of this guide, but shouldn't be too difficult for someone with experience on these platforms.

For either of these commands to work correctly, you will need to make sure that your source directory is laid out in the standard `stele` format.

## File structure

The `stele` CLI relies on a specific file structure to work correctly:

```
.
├── config.yml
├── pages
└── posts
```

### `config.yml`

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
* `menu` - A list of links to include in the top-level site navigation
  ```
  menu:
    - label: about
      path: /about
  ```
  > The `path` should correspond to a static page e.g. `pages/about.html` or the link will be dead when the site is generated.

### `pages/`

"Pages" are static content written in HTML that require no additional parsing. These are primarily useful for things like "About Me" and "Contact" pages that are generally authored once and only infrequently updated.

`stele` will look for any files in the `pages/` subdirectory with an `.html` extension and render them into the default layout. The resulting file will have the same name as the original but will be in the root of the rendered static assets. For example, the file `/pages/about.html` will be available at `{{ BASE_URL }}/about.html`. For this reason, make sure that your file names are URL-safe.

### `posts/`

Posts in `stele`, as with most SSGs, are written in markdown with some minimal frontmatter.

`stele` will look for any files in the `posts/` subdirectory with a `.md` extension and do some processing based on the file name and frontmatter. Posts should be named following the format `YYYYMMDD-post-name-here.md` where `YYYYMMDD` is the "authored" date of the post and `post-name-here` is a descriptive name for the post. As with pages, the post name should be URL-safe as it will be converted in to a slug for the post in the final assets.

The post content should follow the format:

```markdown
---
title: Your Title Here
description: Your description here.
tags: [tag1, tag2]
draft: false
---

Your content here.
```

The `title` and `description` fields are required. If `tags` is not present, the post will simply not have any tags assigned. Finally, if `draft` is `true`, the post will be rendered in the local development server but will not be present in the production build.

---

That's really all you need to get started. Everything else is handled by the framework!

## Alternatives

- [Gatsby](https://www.gatsbyjs.com/)
- [Hugo](https://gohugo.io/)
- [Jekyll](https://jekyllrb.com/)
