# HardhatKV

A Redis-compatible key-value server written in Go. HardhatKV is a fork of [godis](https://github.com/HDT3213/godis).

The documentation site is in [`docs/`](docs/). Pushes to `master` publish it to GitHub Pages.

It listens on `0.0.0.0:6399`. Any Redis client can talk to it, including `redis-cli`.

## Features

- String, list, hash, set, sorted set, and bitmap
- Concurrent core
- TTL
- Publish/subscribe
- GEO
- AOF and AOF rewrite
- RDB read and write
- Multiple databases and `SELECT`
- Transactions are atomic and isolated. If a command in the transaction fails, HardhatKV rolls back the commands that already ran
- Replication
- A server-side cluster that is transparent to the client. A connection to any node can read and write every key in the cluster
  - Raft keeps the cluster metadata. Nodes can join, rebalance, and fail over
  - `MSET`, `MSETNX`, `DEL`, `Rename`, and `RenameNX` run atomically in cluster mode, including when the keys live on different nodes
  - `MULTI` transactions run inside one slot in cluster mode

## Get started

From this directory:

```bash
go run .
```

Or build a binary:

```bash
go build -o hardhatkv .
./hardhatkv
```

`build-darwin.sh`, `build-linux.sh`, and `build-all.sh` write platform binaries under `target/`.

HardhatKV reads the config path from the `CONFIG` environment variable. If that is unset, it reads `redis.conf` in the working directory. Every option is listed in [example.conf](./example.conf).

### Cluster mode

`node1.conf` and `node2.conf` start a two-node cluster:

```bash
CONFIG=node1.conf ./hardhatkv &
CONFIG=node2.conf ./hardhatkv &
```

Connect to either node to reach every key:

```bash
redis-cli -p 6399
```

Cluster options are in [example.conf](./example.conf).

## Supported commands

See [commands.md](./commands.md).

## Benchmark

Environment:

Go version: 1.23
System: MacOS Monterey 12.5 M2 Air

`redis-benchmark` results from the upstream godis tree:

```
PING_INLINE: 179211.45 requests per second, p50=1.031 msec
PING_MBULK: 173611.12 requests per second, p50=1.071 msec
SET: 158478.61 requests per second, p50=1.535 msec
GET: 156985.86 requests per second, p50=1.127 msec
INCR: 164473.69 requests per second, p50=1.063 msec
LPUSH: 151285.92 requests per second, p50=1.079 msec
RPUSH: 176678.45 requests per second, p50=1.023 msec
LPOP: 177619.89 requests per second, p50=1.039 msec
RPOP: 172413.80 requests per second, p50=1.039 msec
SADD: 159489.64 requests per second, p50=1.047 msec
HSET: 175131.36 requests per second, p50=1.031 msec
SPOP: 170648.45 requests per second, p50=1.031 msec
ZADD: 165289.25 requests per second, p50=1.039 msec
ZPOPMIN: 185528.77 requests per second, p50=0.999 msec
LPUSH (needed to benchmark LRANGE): 172117.05 requests per second, p50=1.055 msec
LRANGE_100 (first 100 elements): 46511.62 requests per second, p50=4.063 msec
LRANGE_300 (first 300 elements): 21217.91 requests per second, p50=9.311 msec
LRANGE_500 (first 500 elements): 13331.56 requests per second, p50=14.407 msec
LRANGE_600 (first 600 elements): 11153.25 requests per second, p50=17.007 msec
MSET (10 keys): 88417.33 requests per second, p50=3.687 msec
```

## Layout

- project root: the process entry point
- config: config parser
- interface: interface definitions
- lib: logger, sync helpers, and wildcard matching
- tcp: the TCP server
- redis: the Redis protocol parser
- datastruct: data structures
  - dict: a concurrent hash map
  - list: a linked list
  - lock: key locks
  - set: a hash set
  - sortedset: a skiplist sorted set
- database: the storage engine
  - server.go: a standalone server with multiple databases
  - database.go: one database
  - exec.go: command dispatch
  - router.go: the command table
  - keys.go, string.go, list.go, hash.go, set.go, sortedset.go: command handlers
  - pubsub.go: publish and subscribe
  - aof.go: AOF persistence and rewrite
  - geo.go: geography commands
  - sys.go: authentication and other system commands
  - transaction.go: local transactions
- cluster:
  - cluster.go: cluster mode
  - com.go: node-to-node communication
  - del.go, mset.go, rename.go: atomic multi-key commands
  - keys.go: key commands
  - multi.go: distributed transactions
  - pubsub.go: publish and subscribe in the cluster
  - tcc.go: try-commit-catch transactions
- aof: AOF persistence

## License

HardhatKV is licensed under GPL-3.0. See [LICENSE](./LICENSE).
