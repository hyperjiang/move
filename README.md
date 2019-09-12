# move

[![Docker Repository on Quay](https://quay.io/repository/hyper/move/status "Docker Repository on Quay")](https://quay.io/repository/hyper/move)
[![License](https://img.shields.io/github/license/hyperjiang/move.svg)](https://github.com/hyperjiang/move)

Docker image for dumping and loading mysql schema and data, built on alpine.

## usage

Edit `config.toml` and run something like this:

```
docker run --rm -v $PWD/config.toml:/app/config.toml quay.io/hyper/move
```

The default config path in the container is `/app/config.toml`, if you mount to another position, you need to specify the path like this:

```
docker run --rm -v $PWD/config.toml:/etc/config.toml quay.io/hyper/move -c /etc/config.toml
```

The sql files dumped from database is stored in `/data`, if you want to get them you can mount a local volume to it like this:

```
docker run --rm -v $PWD/config.toml:/app/config.toml -v $PWD/data:/data quay.io/hyper/move
```

## notes

- The program does not support dumping all the databases, this is intended, because dumping once for all the databases is easy to cause problem.
