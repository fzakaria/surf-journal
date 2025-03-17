# Surf Journal

A simple website to track your surfing session.

## Development

The primary way to develop this codebase is using [nix](https://nixos.org).

```console
> nix develop
> go build
# or you can build it directly
> nix build
```

I like to use [entr](https://eradman.com/entrproject/) to _hot reload_ my project.

```console
> ls **/*.go **/*.html.tmpl | , entr -r go run cmd/server/main.go
```