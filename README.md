# Maya

[![Build Status](https://travis-ci.org/if1live/maya.svg?branch=master)](https://travis-ci.org/if1live/maya)
[![Coverage Status](https://coveralls.io/repos/github/if1live/maya/badge.svg?branch=master)](https://coveralls.io/github/if1live/maya?branch=master)

Markdown preprocessor for static site generator.

## Feature
### Generate markdown file from markdown template file.
There are many static site generator exists.
Static site generator requires some metadata. (For example, title, slug, category, tags,...)
There is no standard markdown syntax for metadata.
So, every static site generate make their own syntax to express metadata.

For example, [pelican](http://blog.getpelican.com/) use this markdown.

```
Title: My super title
Date: 2010-12-03 10:20
Modified: 2010-12-05 19:30
Category: Python
Tags: pelican, publishing
Slug: my-super-post
Authors: Alexis Metaireau, Conan Doyle
Summary: Short version for index and feeds

This is the content of my super blog post.
```

[Hugo](https://gohugo.io/) use this markdown.

```
+++
date = "2015-01-08T08:36:54-07:00"
draft = true
title = "about"

+++

## A headline

Some Content
```

If syntax to express metadata exists, we can migrate from pelican to hugo easily.
(or migrate from A-static-site-generator to B-static-site-generator)

### Replace code and command line output
Embedding code into markdown is bothering task. Maya read source and embed it into markdown document.
Embedding command line output into markdown is bothering task. Maya execute command and embed result into makrdown document.


## Install

```bash
go install github.com/if1live/maya
```
## Usage

### Step1. Prepare markdown-like file and other file.

**demo.md**

```md
title: this is title
subtitle: this is subtitle
tags: lisp, demo
author: if1live
slug: sample-article

## write article

~~~maya:view
file=demo.lisp
lang=lisp
~~~

~~~maya:view
file=demo.lisp
start_line=0
end_line=1
lang=lisp
~~~

~~~maya:execute
cmd=clisp demo.lisp
~~~

~~~maya:execute
cmd=clisp demo.lisp
format=blockquote
~~~

```

**demo.lisp**
demo.lisp is used in ``demo.md``.

```lisp
(defun function... () (princ "and ..."))
(function...)

```

## Step 2. Build document

```bash
maya -mode=pelican -file=demo.md
```

```
title: this is title
subtitle: this is subtitle
tags: lisp, demo
author: if1live
slug: sample-article

## write article

```lisp
(defun function... () (princ "and ..."))
(function...)

```

```lisp
(defun function... () (princ "and ..."))
```

```
and ...

```

> and ...
>
>

```

Output is markdown syntax, but it is hard to embed markdown document into another document. so, I use blockquote instead of code syntax.

## Is it Useful?

**This `README.md` is generated from `README.tpl.md`.**
**Embedded code and output are generated by maya.**

## Syntax
### Metadata
```
+++
title: this-is-title
subtitle: this-is-subtitle
<key>: <value>
+++
```

### Embed file

ignore first backslash.

```
\~~~maya:view
file=demo.lisp
lang=lisp
start_line=1
end_end=2
format=blockquote
~~~
```

* file: required, file to attach
* lang: optional, language. if not exist, use extension
* start_line: optional, starting line to begin reading include file
* end_line: optional, last line from include file to display
* format: optional, blockquote/code/bold

### Embed command output

ignore first backslash.

```
\~~~maya:execute
cmd=maya -mode=pelican -file=demo.md
format=blockquote
attach_cmd=true
~~~
```

* cmd: required, command to execute
* format: optional, blockquote/code/bold
* attach_cmd: optional, attach cmd or not (if value exist, attach cmd)
