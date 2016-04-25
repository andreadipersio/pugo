Pugo
======

![pugo](https://github.com/andreadipersio/pugo/raw/master/logo.jpg "Pugo Logo")

A github/statical ultra-minimal website generator.

Pugo is intended to be a starting point for busy developers not having
time to mess up with php/node/whatever frameworks.

You're encouraged to fork it and do your own enhancements, customization.

A live example of pugo is my [personal website](http://andreadipersio.com).

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
