# nixos-shield

## Enterprise edition

We build and deploy a paid Enterprise edition to order. Examples of what we can develop and extend for your organization are listed below.

- Central control for firewall policy across NixOS servers, with per-host templates and rule previews.
- Integrations with Git and deployment pipelines, including versioned policy diffs and rollback.
- Separate author, reviewer and deployment roles with an approval log.

For an Enterprise edition or custom development, contact us at [verobots.eu](https://verobots.eu/) or [olarian.com](https://olarian.com/).


**A security guardian for self-hosted federated services on NixOS.**

Status: **draft specification and working reference CLI.** The CLI converts the network portion
of a service security contract into an nftables ruleset. It does not apply the ruleset.

## The idea in one paragraph

Securing a self-hosted Mastodon, Matrix or mail server is hard. You have to lock the server down, but
without cutting the federation traffic that keeps it connected. Block too much and it goes silent; block
too little and it gets hacked; one wrong rule and you lock yourself out. Normal security tools are
generic: they can't tell one service from another, so they miss what is specific to each. nixos-shield
takes a different route. Each service declares, in a small file, how it is meant to behave. nixos-shield
locks everything else down and watches for anything outside that declaration.

That file is the **security contract**, and its format is what this repository specifies.

## What is here

- [`SPEC.md`](SPEC.md) - the service security contract format (draft).
- **A working reference tool** (`main.go`) that reads a contract and generates a real nftables firewall:
  default-deny, with only what the contract declares, and federation kept alive.
- [`examples/`](examples/) - contracts for Mastodon and Matrix.
- [`CONTRIBUTING.md`](CONTRIBUTING.md) - how to comment and help shape the standard.
- `LICENSE` - AGPL-3.0-or-later.

## Quick start

Download the archive for your OS and CPU from [Releases](../../releases) and extract it.
The archive includes the executable, example contracts and the license.

On Linux, open a terminal in the extracted directory and run:

```sh
./nixos-shield examples/mastodon.json > mastodon.nft
cat mastodon.nft
```

On Windows, preview the same rules without applying them:

```powershell
.\nixos-shield.exe examples/mastodon.json
```

To build locally, use Go 1.23 or newer and run `go build .`. Run `go test ./...` to check
contract validation and rule generation. No third-party Go dependencies are required.

On Linux, validate the generated file using `sudo nft --check --file mastodon.nft` in an isolated
test machine with nftables installed. The rules use default-deny on the entire host. Check access
requirements such as SSH and monitoring before deploying a ruleset.

Output for `examples/mastodon.json`:

    table inet nixos_shield {
      chain input {
        type filter hook input priority 0; policy drop;
        ct state established,related accept
        iif "lo" accept
        ip protocol icmp accept
        ip6 nexthdr icmpv6 accept
        tcp dport 443 accept   # inbound (contract)
      }
      chain output {
        type filter hook output priority 0; policy drop;
        ct state established,related accept
        oif "lo" accept
        tcp dport 443 accept   # https
        udp dport 53 accept   # dns
        tcp dport 53 accept   # dns
      }
    }

The contract said the service listens on 443 and needs ActivityPub federation plus DNS. Everything else,
in and out, is denied. Change the contract, and the firewall changes with it.

## What is not here yet

The behavioral watch, the web-attack layer, containment, and the full NixOS modules. Those are built
next. This tool is the first working slice: the contract really does drive the firewall.

## Who we are

[VeroBots](https://verobots.eu), a small team in Romania that works on servers and security. We want
basic security on the internet to be a given, not something you pay extra for.

## License

AGPL-3.0-or-later. The contract format itself is meant as an open standard: anyone may implement it.
