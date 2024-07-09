# inblog

inblog turns your inbox into a blog

## Usage

inblog depends on a handful of environment variables
```sh
INBLOG_EMAIL='johndoe@example.com' # or sometimes your username
INBLOG_PASSWORD='Password123!' # your password (don't worry, it's sent encrypted)
INBLOG_APPROVED_SENDER='my@email.com' # the email that's allowed to send new posts
INBLOG_IMAPSSERVER='mail.example.com:993' # the IMAPS server
INBLOG_MAILBOX='INBLOG' # the name of the mailbox you want to use
INBLOG_NAME='My Cool inBlog!' # the title of your blog
```

Once these are set, you can start sending emails to the correct place. inblog supports markdown via the CommonMark reference implementation in unsafe mode (meaning you can embed HTML into your markdown for more complicated posts).

Then, to build your blog you just:

```console
$ ./inblog
[INFO] Logging in...
[INFO] 5 new messages!
[INFO] Skipped an email not from approved sender
[INFO] Generating HTML...
```

Your blog will be generated in the `public/` directory, and you can serve this folder directly. If you need to go back and edit a post, you can edit the markdown file and run `./inblog` again to rebuild the page.

You can check out [the repository for my personal blog](https://github.com/weebney/blog) to give you an idea what deploying inblog on GitHub pages is like and how to customize inblog yourself.

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
brew install weebney/inblog
```

### Go toolchain

If you already have the Go toolchain installed, you can use `go install` to install inblog directly.

```console
go install github.com/weebney/inblog@latest
```

If it installs but isn't available, ensure `$(go env GOPATH)/bin` is exported to your `$PATH`.

### Building

A makefile is provided, but inblog can be built easily without it and depends only on the Go toolchain.

```console
go build inblog.go
```

### GitHub releases

Binary releases of inblog are available on the [releases page](https://github.com/weebney/inblog/releases/latest). Just download the correct one and place it somewhere on your `$PATH`.
