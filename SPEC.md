# The service security contract (draft)

Status: **draft, request for comments.** Nothing here is final. The point of publishing early is to
shape it with the people who run federated services.

## Why a contract

Most security tools are generic. They watch for known-bad patterns and ban IP addresses after the fact.
They do not know what a given service is *supposed* to do, so two things happen:

- They miss abuse that looks normal for a generic server but is wrong for this one.
- They break the service when a strict rule blocks traffic the service actually needs (federation is the
  classic case).

A **contract** flips this around. The service states, up front, how it is meant to behave. The guardian
then does two jobs: it closes everything the contract does not allow, and it flags anything the running
service does that the contract did not declare.

The contract can be shipped by the app's own maintainers, or provided alongside (we provide vetted ones
for Mastodon and Matrix to start).

## Shape of a contract

A contract is a declarative document. Below it is shown as a NixOS-style attribute set, since that is the
first target, but the format is meant to be portable (a plain TOML/JSON mapping is equivalent).

```nix
{
  service = "mastodon";
  version = "4.x";

  network = {
    # Ports this service listens on. Everything else inbound is closed.
    inbound = [ 443 ];

    # Outbound is closed by default. Federation needs dynamic destinations,
    # so we name the *flow*, not fixed IPs.
    federation = [ "activitypub" ];      # keep server-to-server ActivityPub alive
    outbound   = [ "https" "dns" ];      # what the service legitimately needs
  };

  web = {
    # Requests that must never be accepted, as patterns. If any of these
    # appears in a URL or body, it is an attack, full stop.
    forbidden = [ "*/wp-login.php" "*/.env" "*/.git/*" "*.sql" ];

    # Reuse an existing web-attack ruleset (OWASP CRS via Coraza), tuned so
    # it does not break signed ActivityPub requests.
    waf = "owasp-crs";
  };

  access = {
    # Devices allowed or denied for admin/privileged access.
    devicesAllow = "./devices/allow.txt";
    devicesDeny  = "./devices/deny.txt";

    # forced: every client must present a usable fingerprint (JA3/JA4).
    # none:   this service is reached in a way that has no fingerprint.
    fingerprint = "forced";

    # Admin surface is closed until explicitly opened (zero-trust).
    adminOnDemand = true;
  };

  behavior = {
    # What "normal" looks like. Deviations are flagged.
    spawnsShell = false;              # this service never spawns a shell
    outboundConnectsTo = "federation"; # its outbound is federation only
    maxConnRate = 200;               # connections/sec considered normal
  };

  contain = {
    # What to do when the contract is violated.
    onViolation = "isolate-and-alert"; # alert | isolate | isolate-and-alert
    rollbackOnCrashLoop = true;        # return to last known-good config
  };
}
```

## Field reference

### `network`
- **`inbound`** - the only ports accepted from outside. Everything else is dropped (default-deny).
- **`federation`** - named federation flows to keep alive (`activitypub`, `matrix`, `smtp`). The
  guardian knows these need dynamic peers and does not break them. This is the field that stops a
  hardened server from going silent on the fediverse.
- **`outbound`** - other legitimate outbound needs. Everything not listed is blocked, which stops a
  compromised service from being used for spam, scanning or exfiltration.

### `web`
- **`forbidden`** - patterns that are always an attack for this service. A Mastodon server has no
  `wp-login.php`; a request for it is hostile by definition.
- **`waf`** - which existing web-attack ruleset to apply (we reuse Coraza + OWASP CRS, not reinvent it).

### `access`
- **`devicesAllow` / `devicesDeny`** - device rules for privileged access.
- **`fingerprint`** - `forced` when the service is reached by a browser or client that has a TLS/HTTP2
  fingerprint (JA3/JA4), so bad bots are caught even when they rotate IPs; `none` when there is nothing
  to fingerprint.
- **`adminOnDemand`** - the admin surface stays closed until the owner asks for it (zero-trust access).

### `behavior`
Declared "normal". The behavioral watch compares the running service against this. A Mastodon process
that suddenly spawns a shell, or reaches out to a non-federation address, is doing something its contract
never declared. That is the signal, and it does not need fuzzy guessing.

### `contain`
- **`onViolation`** - alert only, isolate the service, or both.
- **`rollbackOnCrashLoop`** - return to the last known-good configuration, safely, without locking the
  owner out.

## What this is not

- Not general DDoS protection. A single server cannot stop an attack that saturates its uplink.
- It restores *configuration*, not data.
- It does not replace the app's own login or MFA.

## Open questions (for comment)

- Best portable serialization: TOML, JSON, or a Nix module as the source of truth with exports?
- How much behavior can be declared usefully before it becomes brittle?
- How should contracts be versioned and signed, so a service ships a trustworthy one?
- How should this sit next to Fediversity's `simple-nixos-fediverse` and NixOps4?

Comments and issues welcome. This is meant to be shaped together.
