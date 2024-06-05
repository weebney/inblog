# inblog

inblog turns your inbox into a blog

## usage

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

### customization

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


## installation

inblog is dependent on:
- grep, sed, etc.
- curl
- mpack
- cmark

Then, you can just download the script directly:

```console
$ cd blog
$ curl -O https://raw.githubusercontent.com/weebney/inblog/master/inblog
```

or with git

```console
$ git clone https://github.com/weebney/inblog
```
