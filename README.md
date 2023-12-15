# go-lichess #

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/joanlopez/go-lichess/lichess)

go-lichess is a Go client library for accessing the [Lichess.org API (2.0.0)](https://lichess.org/api).

> [!WARNING]
> This library is still under active development, so please be careful if you're considering it for production use.
> Certain aspects like error handling and request/response lifecycle management are still in progress.

## Installation ##

go-lichess is compatible with modern Go releases in module mode, with Go installed:

```bash
go get github.com/joanlopez/go-lichess
```

will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/joanlopez/go-lichess/lichess"
```

and run `go get` without parameters.

Finally, to use the top-of-trunk version of this repo, use the following command:

```bash
go get github.com/joanlopez/go-lichess@main
```

## Usage ##

```go
import "github.com/joanlopez/go-lichess/lichess" // with go modules enabled (GO111MODULE=on or outside GOPATH)
```
Construct a new Lichess client, then use the various services on the client to
access different parts of the Lichess API. For example:

```go
client := lichess.NewClient(nil)

// export games played by user "chucknorris"
games, _, err := client.Games.ExportByUsername(context.Background(), "chucknorris", nil)
```

Some API methods have optional parameters that can be passed. For example:

```go
client := lichess.NewClient(nil)

// export (up to 5) games played by user "chucknorris"
max := 5
games, _, err := client.Games.ExportByUsername(context.Background(), "chucknorris", &lichess.ExportByUsernameOptions{Max: &max})
```

NOTE: Using the [context](https://godoc.org/context) package, one can easily
pass cancellation signals and deadlines to various services of the client for
handling a request. In case there is no context available, then `context.Background()`
can be used as a starting point.

### Authentication ###

Use the `WithAuthToken` method to configure your client to authenticate using an
OAuth token (for example, a [Personal Access Token](https://lichess.org/api#section/Introduction/Authentication)).

```go
client := lichess.NewClient(nil).WithAuthToken("... your access token ...")
```

Note that when using an authenticated Client, all calls made by the client will
include the specified OAuth token. Therefore, authenticated clients should
almost never be shared between different users.