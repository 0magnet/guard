# guard

Checks a Go repository against the naming, interface and init rules of its
`AGENTS.md`, without breaking on what is already there. Each check keeps an
allowlist of today's violations, and the allowlist may only shrink:

- a new violation fails, with the rule it breaks;
- an allowlist entry that no longer matches fails too, so a fixed name has to
  come out of the list;
- an unsorted allowlist fails.

```
guard names    # Data, Info, Manager, bare Get, ...; util/common packages; is/has bools
guard ifaces   # an interface that embeds another interface
guard inits    # func init() outside package main
```

It reads the Go files git tracks (and untracked ones not ignored) in the
working directory, skipping tests, generated files, `vendor`,
`node_modules`, `testdata` and dot directories. It works from the syntax
alone, so it does not catch a banned word inside a longer name such as
`RouteInfo`, and it cannot tell a `Get` an interface requires from one that
was chosen. Those stay with review.

## In a repository

The allowlists live in `scripts/guard/` (`--dir` to change it). A check with no
violations needs no file. To start, record what is there now:

```
go run github.com/0magnet/guard@main names --write
go run github.com/0magnet/guard@main inits --write
```

and add the targets to the Makefile:

```make
GUARD ?= go run github.com/0magnet/guard@main

check-names: ## Fail on a new generic name, generic package or is/has bool (scripts/guard/names-allowlist.txt only shrinks)
	$(GUARD) names

check-inits: ## Fail on a new func init() outside package main (scripts/guard/init-allowlist.txt only shrinks)
	$(GUARD) inits
```

`--list` prints the current violations without checking anything.

The checks come from skycoin/skywire's `scripts/guard`, which carries the same
rules in its `AGENTS.md`.
