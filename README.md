# Surf Journal

> [https://surf-journal.fly.dev/](https://surf-journal.fly.dev/)

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)
![github master branch workflow](https://github.com/fzakaria/surf-journal/actions/workflows/ci.yml/badge.svg?branch=main)

_A ** WIP ** simple website to log surf sessions and view interesting statistics_.

## Development

### Setup
This codebase is primarily built with [Nix](https://nixos.org/) but it should work on non-NixOS machines (it is tested as such on GitHub CI).

If you are using Nix, you should be all set bu running `nix develop` (or using [direnv](https://direnv.net/)).

Otherwise simply run `./bin/setup`

### Deployment

Deployment happens to [Fly.io](https://fly.io/) automatically via GitHub but you can also deploy locally.

Build the container locally.
```sh
nix build ".#container"

# you can test it via the following:
# load the OCI image
docker load < ./result

# run the container!
docker run -p 3000:80 -it -e "RAILS_ENV=development" -e "RAILS_MASTER_KEY=$(cat ./config/master.key)" registry.fly.io/surf-journal:$LATEST_HASH
```

> You can optionally mount a local volume if you want to persist the SQLite database used.