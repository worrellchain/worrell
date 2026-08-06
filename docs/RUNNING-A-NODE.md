# Validator guide: running a Worrell node

This guide describes how to bring up a full **Worrell** node and how to turn it into a validator on the testnet. Worrell is a proof-of-stake blockchain built on the Cosmos SDK.

| Concept | Value |
|----------|-------|
| Binary | `worrelld` |
| Token | WORRELL (base denomination: `uworrell`, 6 decimals) |
| 1 WORRELL | 1,000,000 uworrell |
| Data directory | `~/.worrell` |
| Testnet chain ID | `worrell-testnet-1` |
| Mainnet chain ID | `worrell-1` |
| Min gas price | `0.025uworrell` |

---

## 1. Hardware requirements

### Testnet

| Resource | Recommendation |
|---------|---------------|
| Operating system | Ubuntu 22.04 LTS |
| CPU | 2 vCPU |
| RAM | 4 GB |
| Disk | 100 GB SSD |
| Network | Stable connection, public IP for the P2P ports |

### Mainnet

For a mainnet validator more headroom is recommended, since the node must sign blocks continuously and maintain history:

| Resource | Recommendation |
|---------|---------------|
| Operating system | Ubuntu 22.04 LTS |
| CPU | 4 vCPU or more |
| RAM | 16 GB or more |
| Disk | 500 GB+ NVMe SSD (with room to grow) |
| Network | Stable low-latency connection, dedicated public IP |

> A block is produced every ~5-6 seconds. A slow disk or an unstable network leads to missed blocks and, over time, *jailing* due to *downtime*.

---

## 2. Installation

### Dependencies

- **Go 1.22+**
- **git**, **make**, **build-essential**

Installing dependencies on Ubuntu 22.04:

```bash
sudo apt update && sudo apt install -y git curl build-essential
```

Installing Go (example with Go 1.22):

```bash
curl -LO https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.profile
source ~/.profile
go version
```

### Build from source

```bash
git clone https://github.com/worrellchain/worrell.git
cd worrell
make install
```

`make install` compiles the `worrelld` binary and places it in `$HOME/go/bin`. Alternatively, you can build it with Ignite CLI:

```bash
ignite chain build
```

Verify the installation:

```bash
worrelld version
worrelld --help
```

---

## 3. Node initialization

Initialize the node's local configuration by choosing a *moniker* (the public name of your node):

```bash
worrelld init <moniker> --chain-id worrell-testnet-1
```

This creates the data directory at `~/.worrell` with the following relevant structure:

- `~/.worrell/config/genesis.json` — genesis (replaced in the next step)
- `~/.worrell/config/config.toml` — CometBFT configuration (peers, P2P, RPC)
- `~/.worrell/config/app.toml` — application configuration (gas, API, gRPC)
- `~/.worrell/config/priv_validator_key.json` — validator signing key (protect it!)
- `~/.worrell/config/node_key.json` — node network identity

---

## 4. Joining the testnet

### 4.1 Get the genesis

Replace the locally generated genesis with the network's official genesis, published in
[worrellchain/networks](https://github.com/worrellchain/networks):

```bash
curl -s https://raw.githubusercontent.com/worrellchain/networks/main/worrell-testnet-1/genesis.json \
  -o ~/.worrell/config/genesis.json
worrelld genesis validate-genesis
```

Verify that you have the exact official file by checking its sha256 against the one
published in the [networks README](https://github.com/worrellchain/networks#join-the-testnet):

```bash
shasum -a 256 ~/.worrell/config/genesis.json
# expected: a81c507b12ba0678c3172394ff4bb03e1c3db60050cc5568c127a24ec19378fd
```

### 4.2 Configure peers and seeds

Edit `~/.worrell/config/config.toml` and set `persistent_peers` in the `[p2p]` section
(format `<node_id>@<host>:26656`). The up-to-date peer list is published in
[`worrell-testnet-1/chain.json`](https://github.com/worrellchain/networks/blob/main/worrell-testnet-1/chain.json)
(`peers` section). Current bootstrap peer:

```toml
# ~/.worrell/config/config.toml  -> [p2p]
persistent_peers = "bb9164c1bd9ed9ff2c0fd9e09b23285698e231de@164.68.98.186:26656"
```

### 4.3 Set the minimum gas price

Edit `~/.worrell/config/app.toml` and set the minimum gas price. This is mandatory for the node to accept and process transactions:

```toml
# ~/.worrell/config/app.toml
minimum-gas-prices = "0.025uworrell"
```

### 4.4 Start the node

```bash
worrelld start
```

Running it as a `systemd` service is recommended so that it starts automatically and restarts on failure. Example unit at `/etc/systemd/system/worrelld.service`:

```ini
[Unit]
Description=Worrell node
After=network-online.target

[Service]
User=worrell
ExecStart=/home/worrell/go/bin/worrelld start
Restart=on-failure
RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now worrelld
sudo journalctl -u worrelld -f
```

### 4.5 Wait for synchronization

Before creating a validator, the node must be fully synchronized:

```bash
worrelld status 2>&1 | jq '.sync_info'
```

When `catching_up` is `false`, the node is up to date. The network has **State Sync enabled**, which allows fast synchronization without downloading the entire history.

---

## 5. Create a validator

> You need: a synchronized node (`catching_up: false`) and an account with sufficient WORRELL balance. On testnet you can request funds from the faucet (500 WORRELL per request, rate-limited to once per hour per address):
>
> ```bash
> curl -X POST http://164.68.98.186:4500 \
>   -H "Content-Type: application/json" \
>   -d '{"address":"worrell1YOURADDRESS..."}'
> ```

### 5.1 Create or import a key

```bash
# Create a new key (store the mnemonic phrase in a safe place)
worrelld keys add <key-name>

# Or import an existing one
worrelld keys add <key-name> --recover
```

Check the address and the balance:

```bash
worrelld keys show <key-name> -a
worrelld query bank balances $(worrelld keys show <key-name> -a)
```

### 5.2 Get the consensus public key

```bash
worrelld tendermint show-validator
```

Returns a value in the format `{"@type":"/cosmos.crypto.ed25519.PubKey","key":"..."}` that you will use in `--pubkey`.

### 5.3 The `create-validator` command

The modern Cosmos SDK command takes a JSON file with the validator definition. Create `validator.json`:

```json
{
  "pubkey": {"@type":"/cosmos.crypto.ed25519.PubKey","key":"<YOUR_PUBKEY>"},
  "amount": "20000000000000uworrell",
  "moniker": "<moniker>",
  "identity": "",
  "website": "",
  "security": "",
  "details": "Worrell testnet validator",
  "commission-rate": "0.05",
  "commission-max-rate": "0.25",
  "commission-max-change-rate": "0.01",
  "min-self-delegation": "1000000"
}
```

And submit it:

```bash
worrelld tx staking create-validator validator.json \
  --from <key-name> \
  --chain-id worrell-testnet-1 \
  --gas auto \
  --gas-adjustment 1.5 \
  --gas-prices 0.025uworrell \
  --yes
```

If your binary version uses the flag-based form, the equivalent is:

```bash
worrelld tx staking create-validator \
  --pubkey="$(worrelld tendermint show-validator)" \
  --amount="20000000000000uworrell" \
  --moniker="<moniker>" \
  --commission-rate="0.05" \
  --commission-max-rate="0.25" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1000000" \
  --from=<key-name> \
  --chain-id=worrell-testnet-1 \
  --gas=auto \
  --gas-adjustment=1.5 \
  --gas-prices=0.025uworrell \
  --yes
```

### 5.4 Explanation of the values

| Flag / field | Example value | Meaning |
|--------------|-------------------|-------------|
| `--amount` | `20000000000000uworrell` | Initial self-delegation (in this example, 20,000,000 WORRELL). Adjust to your case. |
| `--commission-rate` | `0.05` | Current commission: 5%. **It cannot be lower than the global network minimum (5%).** |
| `--commission-max-rate` | `0.25` | Maximum commission **this validator** allows itself: 25%. You set it. |
| `--commission-max-change-rate` | `0.01` | Maximum commission change per day: 1%. You set it. |
| `--min-self-delegation` | `1000000` | Minimum self-delegation you will maintain. **`1000000` uworrell = 1 WORRELL.** |

> ⚠️ **Critical — `--min-self-delegation`:** the value is expressed in `uworrell`, not in WORRELL. To guarantee the network minimum of 1 WORRELL you must use **`--min-self-delegation="1000000"`**. **NEVER** use `"1"`: that would mean 1 uworrell (0.000001 WORRELL) and would leave your validator below the intended effective minimum.

### 5.5 Commissions: global minimum vs. your own maximum

- The **`min_commission_rate` of 5% (0.05) is GLOBAL**: it is the minimum imposed by the network. No validator can set a `commission-rate` below 0.05.
- In contrast, the **maximum commission** (`commission-max-rate`) and the **maximum daily change** (`commission-max-change-rate`) are decided by **each validator on its own**. Once set when creating the validator, the `commission-max-rate` **cannot be increased afterward**, so choose it carefully.

### 5.6 Verify that the validator is active

```bash
worrelld query staking validator $(worrelld keys show <key-name> --bech val -a)
```

Check that the status is `BOND_STATUS_BONDED` and that `jailed` is `false`.

---

## 6. Monitoring

### 6.1 Node status

```bash
worrelld status
worrelld status 2>&1 | jq '.sync_info.latest_block_height, .sync_info.catching_up'
```

### 6.2 Uptime and signed blocks

Check your validator's signing information (you need its *consensus address*):

```bash
# Validator consensus address
worrelld query slashing signing-info $(worrelld tendermint show-address)
```

The result includes `missed_blocks_counter` (blocks not signed within the window) and `jailed_until`. Keep an eye on `missed_blocks_counter` staying low.

### 6.3 Relevant slashing parameters (avoiding jailing)

| Parameter | Value | What it means |
|-----------|-------|---------------|
| Signed blocks window | 10,000 blocks | Period over which availability is measured |
| Minimum signed per window | 5% (0.05) | You must sign at least 5% of the blocks in the window |
| Downtime penalty | 0.01% (0.0001) | Slash when jailed for inactivity |
| Double-signing penalty | 5% (0.05) | Slash for signing two blocks at the same height |
| Jail duration for downtime | 300s (5 min) on testnet / 43,200s (12 h) on mainnet | Minimum time jailed |

To **avoid jailing due to downtime**:

- Run the node as a `systemd` service with automatic restart.
- Keep disk and CPU with headroom; a node that cannot keep up misses blocks.
- **Never** run two instances with the same `priv_validator_key.json` at the same time: it would cause a **double sign** and a 5% slash (much more serious than downtime).
- Monitor with Prometheus (see below) and set up alerts on missed blocks.

### 6.4 Recovering a jailed validator (unjail)

If your validator was jailed for downtime, once the jail period has elapsed you can reactivate it:

```bash
worrelld tx slashing unjail \
  --from <key-name> \
  --chain-id worrell-testnet-1 \
  --gas auto \
  --gas-adjustment 1.5 \
  --gas-prices 0.025uworrell \
  --yes
```

### 6.5 Prometheus

Prometheus telemetry is enabled. To expose metrics, in `~/.worrell/config/config.toml`:

```toml
# ~/.worrell/config/config.toml  -> [instrumentation]
prometheus = true
prometheus_listen_addr = ":26660"
```

Metrics become available at `http://<host>:26660/metrics`.

---

## 7. Ports

Make sure the firewall allows the necessary traffic. The P2P port must be reachable from outside; the rest should preferably be restricted to localhost or trusted IPs.

| Service | Port | Recommended exposure |
|----------|--------|------------------------|
| CometBFT P2P | 26656 | Public (required to connect to the network) |
| CometBFT RPC | 26657 | Localhost / trusted IPs |
| REST API | 1317 | Localhost by default |
| gRPC | 9090 | Localhost by default |
| Prometheus | 26660 | Localhost / internal monitoring network |

By default, the REST API and gRPC listen only on localhost. If you expose them publicly, protect them with a reverse proxy and firewall rules.

---

## 8. Quick checklist

- [ ] Go 1.22+ installed and `worrelld version` works
- [ ] Node initialized with `--chain-id worrell-testnet-1`
- [ ] Official `genesis.json` placed and validated
- [ ] `seeds` / `persistent_peers` configured
- [ ] `minimum-gas-prices = "0.025uworrell"` in `app.toml`
- [ ] Node synchronized (`catching_up: false`)
- [ ] Account with sufficient balance
- [ ] Validator created with `--min-self-delegation="1000000"` (NOT `"1"`)
- [ ] `commission-rate` ≥ 0.05 (global network minimum)
- [ ] Node running under `systemd` with automatic restart
- [ ] Signed-blocks monitoring active
