# lastgo

The **lastgo** project provides an alternative `go` command which is a proxy keeping up to date the [Go](https://go.dev) toolchain.

## Getting started

Install via [Homebrew](https://brew.sh)

```console
$ brew tap dvaumoron/tap
$ brew install lastgo
```

Or get the [last binary](https://github.com/dvaumoron/lastgo/releases) depending on your OS.

## Environment Variables

### LASTGO_ASK

String (Default: "")

When non empty, **lastgo** ask user for confirmation before an update.

### LASTGO_CHECK_INTERVAL

String (Default: 24h)

Minimum interval waited between check of last version on `LASTGO_DOWNLOAD_URL`.

### LASTGO_DOWNLOAD_URL

String (Default: https://go.dev/dl)

URL to download Go binary (only used to install at first launch, update rely on **Go** `GOTOOLCHAIN` mecanism).

### LASTGO_ROOT

String (Default: ${HOME}/.lastgo)

Path to directory where **lastgo** install `go`.
