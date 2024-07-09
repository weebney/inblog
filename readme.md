# inblog

inblog turns your inbox into a blog

## Usage

inblog is a static site generator (e.g. Hugo) that logs into an email account, downloads the emails, and then converts them into a blog.

### Quickstart

Running `inblog` for the first time, you'll be thrown into a setup wizard which will allow you to configure inblog for use on your machine. Once everything is set up, you can start sending emails to the correct place. Then, to build your blog you just:

```console
$ inblog
```

Your blog will be statically generated in the `./public/` directory, and you can serve this folder directly. If you need to go back and edit a post, you can edit the markdown file and run `inblog -r` to rebuild the page.

### Advanced Usage

In the interest of making it possible to automate inblog in a public GitHub Actions runner, you can also configure inblog entirely via environment variables:

```sh
INBLOG_EMAIL='inblogEmail@example.com' # your inblog email
INBLOG_PASSWORD='Password123!' # your password
INBLOG_APPROVED_SENDER='my@personalemail.com' # the email that's allowed to send new posts
INBLOG_IMAPSSERVER='mail.example.com:993' # the IMAPS server
INBLOG_MAILBOX='INBOX' # the name of the mailbox you want to use
INBLOG_NAME='My inblog' # the title of your blog
```

You can check out [the repository for my personal blog](https://github.com/weebney/blog) to give you an idea what deploying inblog on GitHub pages is like and how to customize inblog yourself.

There are also a handful of command line options that can be passed:

```console
$ inblog -h
Usage: inblog [OPTIONS...]
  -c, --content-dir string   set the content directory (default "content")
  -h, --help                 show this help message
  -m, --no-wizard            disable the setup wizard
  -o, --output-dir string    set the output directory (default "public")
  -r, --rebuild              force a rebuild of the HTML
  -R, --refresh              refresh all content, overwriting edits
  -s, --skip                 skip fetching emails
  -v, --version              print version info and quit
  -w, --wizard               force the setup wizard
inblog v1.0 (c) weebney 2024, BSD-2-Clause
```

-----

### Customization

You can fully customize the look of your inblog by editing the HTML templates found in the `content/templates/` directory. You can utilize special tokens that will expand to various pieces of content relevant to these pages.

The following tokens are avaliable in `post.template.html` and `listitem.template.html`:

- `%SUBJECT%`
- `%DATE%`
- `%BLOG_NAME%`
- `%BODY%`

The following tokens are avaliable for `listitem.template.html`:

- `%HYPERLINK%`

The following tokens are avaliable for `index.template.html`:

- `%BLOG_NAME%`
- `%LIST%`, which is just every post's `listitem.template.html`

-----

## Installation

### Homebrew

inblog is available as a Homebrew package on Linux and macOS

```console
$ brew install weebney/tap/inblog
```

### Go toolchain

If you already have the Go toolchain installed, you can use `go install` to install inblog directly.

```console
$ go install github.com/weebney/inblog@latest
```

If it installs but isn't available, ensure `$(go env GOPATH)/bin` is exported to your `$PATH`.

### Building from source

A makefile is provided, but inblog can be built easily without it and depends only on the Go toolchain.

```console
$ go build inblog.go
```

### GitHub releases

Binary releases of inblog are available on the [releases page](https://github.com/weebney/inblog/releases/latest). Just download the correct one and place it somewhere on your `$PATH`.
