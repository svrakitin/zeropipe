# zeropipe

`zeropipe` is a pipe over the local network using [zeroconf](https://github.com/grandcat/zeroconf) for receiver discovery.

# Usage

`zeropipe` fanouts standard input over TCP to each receiver standard output on cooldown. This cooldown resets every time there is a new receiver registered under `{pipe_id}._zeropipe._tcp`. Cooldown can be configured through `-c/--cooldown` flag or `ZEROPIPE_COOLDOWN` env variable.

Sender:
```sh
# Sender
$ echo "Hello" | zeropipe send my-pipe-id -c 5s
```

Receiver:
```sh
$ zeropipe recv my-pipe-id
Hello
```

If you want to authenticate (challenge) sender on receiver side you can use challenge token (`-t` flag or `ZEROPIPE_TOKEN` env variable).

```sh
$ ZEROPIPE_TOKEN=secret echo "Hello" | zeropipe send my-pipe-id -c 5s
```
```sh
$ zeropipe recv my-pipe-id -t secret
Hello
```
