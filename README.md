# move

[![Docker Repository on Quay](https://quay.io/repository/hyper/move/status "Docker Repository on Quay")](https://quay.io/repository/hyper/move)
[![Go Report](https://goreportcard.com/badge/github.com/hyperjiang/move)](https://goreportcard.com/report/github.com/hyperjiang/move)
[![License](https://img.shields.io/github/license/hyperjiang/move.svg)](https://github.com/hyperjiang/move)

Docker image for dumping and loading mysql schema and data, built on alpine.

The original intention of this tool is to solve this issue: need to export the entire database structure, and only include certain tables' data.

## usage

You can take a look at `config.toml` to understand the features this tool supports, and create your own config and run command like this:

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

If you only want to run one of the rules defined in your `config.toml`, and its name is `r1`, then you can run it by this:

```
docker run --rm -v $PWD/config.toml:/app/config.toml quay.io/hyper/move -r r1
```

## notes

- The program does not support dumping all the databases in one time, this is intended, you need to specify the database in each rule.
- You can run the program directly without docker if you are using linux or mac and have mysql client installed on your local.
- The rules are run in FIFO order.
