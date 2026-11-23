# Contributing

This is an early draft of an open standard. The most useful thing right now is your judgement, not code.

If you run a federated service, or maintain one, we want to know:

- Does the contract capture how your service actually behaves?
- What did we get wrong, or leave out?
- Where would this clash with how you already deploy (reverse proxy, firewall, NixOps4, SelfHostBlocks)?

Open an issue, or send comments. Rough notes are fine.

## Ground rules

- The contract format stays simple and readable. If a field needs a paragraph to explain, it probably
  does not belong.
- We reuse existing engines (Coraza, standard scanners) rather than reinvent them.
- Everything here is AGPL-3.0-or-later. The format itself is open for anyone to implement.
