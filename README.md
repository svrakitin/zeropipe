# zeropipe

`zeropipe` is a pipe over the local network using [zeroconf](https://github.com/grandcat/zeroconf) for receiver discovery.

# Usage

`zeropipe` fanouts standard input over TCP to each receiver standard output on cooldown. This cooldown resets every time there is a new receiver registered under `{pipe_id}._zeropipe._tcp`. Cooldown can be configured through `-c/--cooldown` flag or `ZEROPIPE_COOLDOWN` env variable.

Sender:
```sh
$ echo "Hello" | zeropipe send my-pipe-id -c 5s
```

Receiver:
```sh
$ zeropipe recv my-pipe-id
Hello
```
