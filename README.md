# smon

`smon` is a service monitor that I wrote to debug what's wrong with my home
network. It's potentially usable for a broader application of system
monitoring.

## Invoking smon & checkers

The invocation is very simple:

```shell
smon CONFIG.json
```

starts service checkers that `CONFIG.json` requests. At the time of writing
the following checkers exist:

* Is there a local IP address (other than loopback)? If no, the WiFi may be
   down, its DHCP may be unreachable or not working.
* Can DNS resolve a hostname? May point to a DNS error.
* Can a webpage be retrieved?
* Can a route to a host be resolved (`traceroute`)?

More checkers might have been added; seen `smon.json.sample`.

## Configuration

The configuration is fairly simple: it's an array of JSON objects, where each
object states:

* The `checker`, which must be an identifier that `smon.go` can recognize
* A one-string argument `arg`, if the checker needs one. E.g., the host lookup
  checker requires one. The `arg` is not given for checkers that don't need it,
  such as the verification that a local IP is known.
* An `interval` stating how often the checker should be run. The format is
  e.g. `10s` for each 10 seconds, `5m` for each 5 minutes and so on.
* Optionally a `maxduration` to limit the checker runtime. When not given, the
  maximum runtime defaults to the `interval`.
