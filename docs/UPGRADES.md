# Chain upgrades

This guide explains how coordinated software upgrades work on Worrell and gives the
step-by-step procedure for each upgrade.

| Upgrade | From | Upgrade height | Release |
|---------|------|----------------|---------|
| [`v0.1.3`](#upgrade-v013) | `v0.1.2` | To be announced | [v0.1.3](https://github.com/worrellchain/worrell/releases/tag/v0.1.3) |

> **Stay in the loop.** Upgrade heights and reminders are published on Telegram:
> [t.me/worrellvalidators](https://t.me/worrellvalidators).

---

## How a coordinated upgrade works

1. A new `worrelld` release is published. It contains an *upgrade handler* with a name (for example `v0.1.3`).
2. A governance **software-upgrade proposal** schedules that name at a block height.
3. If the proposal passes, every node running the old binary **stops producing blocks at that height** and logs `UPGRADE "<name>" NEEDED at height: <height>`.
4. Each validator replaces the binary and starts the node again. The new binary runs the upgrade handler in that same block.
5. The chain resumes as soon as validators holding more than 2/3 of the voting power are running the new binary.

Nothing has to be done at an exact second: the old binary always stops by itself at the
upgrade height. What matters is replacing the binary **after** the node has stopped, and
doing it soon so the chain can resume.

---

## Upgrade v0.1.3

| Field | Value |
|-------|-------|
| Upgrade name | `v0.1.3` |
| Upgrade height | **To be announced** in the governance proposal and on Telegram |
| Upgrades from | `v0.1.2` |
| Release | [v0.1.3](https://github.com/worrellchain/worrell/releases/tag/v0.1.3) (commit `a914f444004df7514ee909c4fa2e66a942059ca6`) |
| Cosmos SDK | v0.53.6 (unchanged) |
| State changes | None |

**What changes.** This is the first coordinated upgrade of the testnet and its purpose is
to exercise the upgrade procedure. The new binary registers the `v0.1.3` upgrade handler,
which runs the pending module migrations; no module changes version, so balances,
validators, delegations and parameters are left untouched. Release binaries are now built
with a pinned Go toolchain (Go 1.26.5).

### Release checksums (sha256)

| File | sha256 |
|------|--------|
| `v0.1.3_linux_amd64.tar.gz` | `b6e8687d1af51f8dea81c1ba6ecd93b358ce6bb80468159c1c7f70d4b908a6ec` |
| `v0.1.3_linux_arm64.tar.gz` | `11c2e324b03693380515e97b8090a86f132bc4ea467d80d0e8d05709d15ee5d0` |
| `v0.1.3_darwin_amd64.tar.gz` | `a07c18411767779195f20b9d911459149f56e6b552624b4b964442c05a9d63b5` |
| `v0.1.3_darwin_arm64.tar.gz` | `6bc33804e98058443568dabcd77330940e39007224a2b874116cff58c6f4cf8c` |

### 1. Before the upgrade height

> ⚠️ **Do not switch to the new binary before the node has stopped at the upgrade height.**
> Once the proposal has passed, a node started with `v0.1.3` before that height refuses to
> continue with `BINARY UPDATED BEFORE TRIGGER! UPGRADE "v0.1.3" - in binary but not executed on chain`.
> If that happens, put the `v0.1.2` binary back and restart: the node rejoins normally.

Download the new binary, verify it and leave it ready in a separate directory:

```bash
cd ~
curl -LO https://github.com/worrellchain/worrell/releases/download/v0.1.3/v0.1.3_linux_amd64.tar.gz
curl -LO https://github.com/worrellchain/worrell/releases/download/v0.1.3/release_checksum
sha256sum -c release_checksum --ignore-missing     # must print: v0.1.3_linux_amd64.tar.gz: OK
mkdir -p ~/worrell-v0.1.3 && tar -xzf v0.1.3_linux_amd64.tar.gz -C ~/worrell-v0.1.3
~/worrell-v0.1.3/worrelld version --long | grep -E '^(version|commit|go):'
```

Expected output:

```
commit: a914f444004df7514ee909c4fa2e66a942059ca6
go: go version go1.26.5 linux/amd64
version: 0.1.3
```

Also, ahead of time:

- **Check your system.** The `linux_amd64` binary needs glibc 2.34 or newer (Ubuntu 22.04+,
  Debian 12+); check with `ldd --version`. On older systems it fails to start with
  `version 'GLIBC_2.34' not found`: [build from source](#4-building-from-source) instead.
  The `linux_arm64` binary is statically linked and has no such requirement.
- **Locate the binary your node is running**, so you replace the right file:

  ```bash
  readlink -f /proc/$(pgrep -x worrelld | head -1)/exe
  ```

- **Back up** `~/.worrell/config/priv_validator_key.json` and `~/.worrell/config/node_key.json`
  if you have not done so already.
- **Check your keys with the new binary.** Run `~/worrell-v0.1.3/worrelld keys list`, with the
  same `--keyring-backend` you normally use, and check that your keys are listed. This does not
  start the node. If the list is empty, see [Building from source](#4-building-from-source).

### 2. Option A — Manual upgrade

1. **Wait for the node to stop at the upgrade height.** Follow the log:

   ```bash
   sudo journalctl -u worrelld -f | grep -E 'UPGRADE|CONSENSUS FAILURE'
   ```

   At the upgrade height it prints `UPGRADE "v0.1.3" NEEDED at height: <height>`.
   **The process does not exit**: the node stays up but no longer produces or signs blocks,
   so `systemd` will not restart it for you.

2. **Stop the node.**

   ```bash
   sudo systemctl stop worrelld
   ```

3. **Replace the binary**, keeping a copy of the old one. Use the path you located before
   (`/home/worrell/go/bin/worrelld` in the example unit of the node guide, `~/bin/worrelld`
   if you installed the prebuilt binary):

   ```bash
   BIN=/home/worrell/go/bin/worrelld
   cp "$BIN" "$BIN.v0.1.2"
   cp ~/worrell-v0.1.3/worrelld "$BIN"
   "$BIN" version        # 0.1.3
   ```

4. **Start the node.**

   ```bash
   sudo systemctl start worrelld
   sudo journalctl -u worrelld -f
   ```

   The log shows `applying upgrade "v0.1.3" at height: <height>`. The node then waits until
   more than 2/3 of the voting power has upgraded; a height that does not move for a while
   right after the upgrade is expected.

If the node is restarted with the old binary after the upgrade height, it exits with
`UPGRADE "v0.1.3" NEEDED`. With `Restart=on-failure` it keeps restarting until the binary
is replaced; this is harmless.

### 3. Option B — Cosmovisor

If your node already runs under [Cosmovisor](https://docs.cosmos.network/main/tooling/cosmovisor),
place the new binary before the upgrade height and Cosmovisor does the switch by itself:

```bash
export DAEMON_NAME=worrelld
export DAEMON_HOME=$HOME/.worrell
cosmovisor add-upgrade v0.1.3 ~/worrell-v0.1.3/worrelld
ls $DAEMON_HOME/cosmovisor/upgrades/v0.1.3/bin/     # worrelld
```

At the upgrade height Cosmovisor stops the old binary, backs up the data directory, switches
to `upgrades/v0.1.3/bin/worrelld` and restarts the node. Keep in mind:

- **The backup needs free disk space and time** proportional to the size of `~/.worrell/data`.
  It can be disabled with `UNSAFE_SKIP_BACKUP=true`, at your own risk.
- **Place the binary yourself** rather than relying on automatic download
  (`DAEMON_ALLOW_DOWNLOAD_BINARIES=false`), so you run the binary you verified.

### 4. Building from source

Build from the **release tag**, not from `main`, and use **Go 1.26.5** to match the official
binaries (`go.mod` accepts older versions, but all validators should run binaries built with
the same compiler):

```bash
cd worrell
git fetch --tags
git checkout v0.1.3
make install          # installs to $HOME/go/bin/worrelld
worrelld version --long | grep -E '^(version|commit|go):'
```

A binary built with `make install` reports `version: v0.1.3` and `name: worrell`, while the
official one reports `version: 0.1.3` and `name: Worrell`. The difference is cosmetic for the
node, with one exception: with `--keyring-backend os` or `pass`, keys are stored under the
binary's name, so keys created with one are not listed by the other. The `file` and `test`
backends, and the validator signing key, are not affected.

Keys that are not listed are not lost: they are still stored under the other name. Either
install `v0.1.3` the same way you installed `v0.1.2`, or export each key with the old binary
(`worrelld keys export <key-name>`) and import it with the new one
(`worrelld keys import <key-name> <file>`), or restore it from its mnemonic
(`worrelld keys add <key-name> --recover`).

### 5. After the upgrade

```bash
worrelld version                                   # 0.1.3
worrelld query upgrade applied v0.1.3              # prints the upgrade height
worrelld status 2>&1 | jq '.sync_info.latest_block_height, .sync_info.catching_up'
worrelld query slashing signing-info $(worrelld tendermint show-address)
```

The height must keep increasing and `missed_blocks_counter` must stop growing.

### 6. Troubleshooting

| Message or symptom | Meaning | What to do |
|--------------------|---------|------------|
| `UPGRADE "v0.1.3" NEEDED at height: <height>` | The node runs `v0.1.2` at or after the upgrade height | Replace the binary with `v0.1.3` and start the node |
| `BINARY UPDATED BEFORE TRIGGER! UPGRADE "v0.1.3"` | The node runs `v0.1.3` before the upgrade height | Put `v0.1.2` back and restart; switch again after the node stops at the height |
| `version 'GLIBC_2.34' not found` | The operating system is too old for the `linux_amd64` binary | Build from source |
| The node runs `v0.1.3` but the height does not move | Less than 2/3 of the voting power has upgraded yet | Wait; nothing to do on your side |
| `worrelld keys list` is empty with the new binary | Keys are stored under the other binary name (`os` or `pass` keyring) | See [Building from source](#4-building-from-source) |

**Upgrading late.** A validator that misses the upgrade height can upgrade afterwards with
the same steps: the node catches up and signs again. Late upgrades delay everyone: the chain
does not resume until more than 2/3 of the voting power has upgraded. Once it resumes, a
validator that is still on the old binary is jailed for downtime after missing about 9,500 blocks.

**New nodes after the upgrade.** A node syncing from genesis must start with `v0.1.2`; it
stops at the upgrade height like everyone else and continues with `v0.1.3`. A node that uses
State Sync from a snapshot taken after the upgrade height starts directly with `v0.1.3`.
