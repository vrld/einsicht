# Einsicht, the mail viewer

Einsicht is a simple tool to view a single mail file. It displays plain text
and opens sanitized html mails in your browser. It previews and saves
attachments. Use it as a TUI or call it from scripts.

## TUI

Open the tui with `einsicht`. The mail file is read from stdin:

    einsicht < /path/to/mail.eml

You can also specify the path with `--file`, `-f`:

    einsicht -f /path/to/mail.eml


## CLI

All commands require input on stdin or the `--file` flag.

List email parts:

    einsicht list parts
    einsicht l p
    einsicht l

    einsicht l bodies
    einsicht l attachments

Print email parts to stdout:

    einsicht print  # prints text/plain body
    einsicht p text
    einsicht p html
    einsicht p attachment <n>

Open email parts; uses `xdg-open` by default, but can be overridden:

    einsicht open
    einsicht o text --command cat
    einsicht o html
    einsicht o att <n>

Save email parts; will open a file dialog via desktop-portal:

    einsicht save
    einsicht s t
    einsicht s h
    einsicht s a <n>

## HTML content

HTML content will be modified:

- prevent loading of external resources (images, styles, etc), unles user allows it
- prevent clicking on links, unless user allows it
