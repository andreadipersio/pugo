Pugo
======

![pugo](https://github.com/andreadipersio/pugo/raw/master/logo.jpg "Pugo Logo")

A github/statical ultra-minimal website generator.

Pugo is intended to be a starting point for busy developers not having
time to mess up with php/node/whatever frameworks.

You're encouraged to fork it and do your own enhancements, customization.

A live example of pugo is my [personal website](http://andreadipersio.com).

### Why?
In 2013 I was in full 'Go' mode, most of my side project and a lot of my daily work involved Go.
Being also an avid fan of Markdown and Github I decided to create a static website generator 
that I could use to present Markdown files in repositories as a HTMl web page that I would use in my 
personal website.

Eventually, spending most of my time working on work related projects I neglected my personal website 
and eventually decided to remove all of its content.

For posterity this is how `https://andreadipersio.com/ds2key-srv` would have looked like.

[ds2key-srv-golang.pdf](https://github.com/andreadipersio/pugo/files/12879787/ds2key-srv-golang.pdf)

### Usage

```shell
./pugo --repo repositoryname --owner ownername --token githubtoken
```

### Urls example:

Suppose you have the following repo:

```
myrepo
|---my-nice-article.md
|---another-nice-article.md
|---README.md
```

You get the following urls:

```
/                       -> README.md
/my-nice-article        -> my-nice-article.md
/another-nice-article   -> another-nice-article.md
```

Once you access a page, it get cached.
To force a cache refresh use the `refresh=True` GET query.

```
/my-nice-article?refresh=True
```
