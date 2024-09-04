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

Assuming `stele` is already installed and in your path, using the CLI is pretty
straightforward. Running

```
stele dev
```

will start up a local development server. This does a full build of the static assets and makes them available at `http://localhost:8081`. At the moment, the site will not automatically rebuild on changes, but you can visit `http://localhost:8081/refresh` to force a rebuild.

Once you are satisfied with your post/page/whatever, you can do a full build by running

```
stele build
```
.

This will create a `dist/` folder with all of the static assets for the site. These assets should be deployable as-is to something like an S3 bucket or you can have your favorite host (e.g. Cloudflare Pages, Netlify, etc.) build and deploy them for you. How to set that up is outside of the scope of this guide, but shouldn't be too difficult for someone with experience on these platforms.

For either of these commands to work correctly, you will need to make sure that your source directory is laid out in the standard `stele` format.

### File structure

The `stele` CLI relies on a specific file structure to work correctly:

```
.
├── config.yml
├── pages
└── posts
```

### `config.yml`

At the root of your project you will need a YAML file containing the configuration for your blog. The available configuration values are:

|name|description|example|
|----|-----------|-------|
|`title`|The name of your blog|`title: My First Blog`|
|`description`|A short description of your blog's content/purpose|`description: A personal blog about fish`|
|`author`|Probably your name|`author: Grace Hopper`|
|`baseURL`|The FQDN and protocol for your blog|`baseURL: https://myblog.example.com`|
|`categories`|A list of categories that describe the content of the blog|`categories: [blog, programming]`|
|`menu`|A list of links to include in the top-level site navigation|`menu:
  - label: about
    path: /about`|

### `pages/`

### `posts/`

## Alternatives

- [Gatsby](https://www.gatsbyjs.com/)
- [Hugo](https://gohugo.io/)
- [Jekyll](https://jekyllrb.com/)
