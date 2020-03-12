# smon

`smon` is a service monitor that I wrote to debug what's wrong with my home
network. It's potentially usable for a broader application of system
monitoring.

## Invoking smon & checkers

The invocation is very simple. In its most basic form, just run:

```shell
smon CONFIG.json
```

This starts the service checkers that `CONFIG.json` requests. You can try
`smon --help` to see relevant flags. For example,

```shell
smon -p 8080 CONFIG.json
```

will fire up a simple HTTP server on port 8080 where you can point your
browser to see some results.

At the time of writing the following checkers exist:

* Is there a local IP address (other than loopback)? If no, the WiFi may be
   down, its DHCP may be unreachable or not working.
* Can DNS resolve a hostname? May point to a DNS error.
* Can a webpage be retrieved?
* Can a route to a host be resolved (`traceroute`)?

See `smon.json.sample` for an example. The JSON configuration must contain one
top-level array `[...]` of JSON objects `{...}`. Every object configures a
checker:

* There must be a key `checker`, stating the name of the checker to run.
* There must be an `interval`; the checker will be run each interval once.
* There may be an `arg` when the checker requires additional information.
* There may be a `maxruntime` stating how long the checker is allowed to run
  before being considered a failure. When absent, `interval` is the used.

Durations (`interval` and `maxruntime`) are expressed as Go durations; e.g.,
`5s` (5 sec), `1h`, `1m30s` etc.

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

## What do the checkers do?

* Checker `WGet` contacts a web URL and fetches a page. It succeeds when the
  response is `200 OK`. It needs a configuration `arg`, the URL.
* Checker `LookupHost` verifies that DNS resolving works. It needs a
  configuration `arg`, a hostname.
* Checker `TraceRoute` runs traceroute on its `arg`.
* Checker `LocalIP` has no `arg` and succeeds when the local system has a
  non-public IP address. This can be used to check that your local DHCP works.

Checker `Dummy` also exists, but it's just for testing. Its `arg` needs 3
numerical values; the chance that it should simulate a failure, the minimum
runtime, and the maximum runtime. See e.g., `smon.json.dummy` for an example.
This causes 25% of reported failures and each checker run lasts between 2 and
5 seconds. However, as the `interval` is set to 4 seconds, that means that
about 1 in 4 checker runs will overrun the interval and will be killed.

## Future work / possible extensions

Checkers are now independent; ie., if you schedule them, then they run.
Expressing dependencies would be nice; e.g., `LookupHost` and `LocalIP`
are useless when `WGet` succeeds. I'm thinking of changing the configuration
and the code to express such a dependency in a *run-this-if-that-fails*
fashion.
