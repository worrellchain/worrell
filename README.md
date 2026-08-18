# Worrell

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://github.com/worrellchain/worrell/blob/main/LICENSE)

**Worrell** is a proof-of-stake (PoS) blockchain built with [Ignite CLI](https://ignite.com/) and [Cosmos SDK](https://docs.cosmos.network/), focused on **payments and energy infrastructure**.

| Field | Value |
|-------|-------|
| Token | WORRELL |
| Base denomination | uworrell |
| Decimals | 6 (1 WORRELL = 1,000,000 uworrell) |
| Binary | `worrelld` |
| Data directory | `~/.worrell/` |
| Address prefix | `worrell` (addresses `worrell1...`) |
| Mainnet chain ID | `worrell-1` |
| Testnet chain ID | `worrell-testnet-1` |
| License | Apache 2.0 |

---

## Key features

- **Proof-of-Stake consensus** over CometBFT, with block times of ~5-6 seconds.
- **Focus on payments and energy infrastructure** as priority use cases.
- **On-chain governance** with standard and expedited (urgent) proposals.
- **Staking and delegation** with up to 100 active validators.
- **Dynamic inflation** tied to the staking ratio (target 67%).
- **Linear, continuous founder vesting**, with multisig accounts for the treasury and community reserves.
- **IBC installed** (disabled at genesis, to be enabled via governance once the network is stable).

---

## Token economics

**Total supply: 1,000,000,000 WORRELL** (1,000,000,000,000,000 uworrell).

Genesis distribution:

```
treasury    30.0%  ██████████████████████████████
airdrop     22.5%  ██████████████████████▌
incentives  17.5%  █████████████████▌
henry       10.0%  ██████████
reserve     10.0%  ██████████
george       5.0%  █████
charlie      5.0%  █████
```

| Account | WORRELL | % | uworrell | Custody |
|--------|--------:|--:|---------:|----------|
| treasury | 300,000,000 | 30.0% | 300000000000000 | Multisig 3/3 |
| airdrop | 225,000,000 | 22.5% | 225000000000000 | Multisig 2/4 |
| incentives | 175,000,000 | 17.5% | 175000000000000 | Multisig 2/4 |
| henry | 100,000,000 | 10.0% | 100000000000000 | Vesting 4 years, cliff 1 year |
| reserve | 100,000,000 | 10.0% | 100000000000000 | Multisig 2/4 |
| george | 50,000,000 | 5.0% | 50000000000000 | Vesting 4 years, cliff 1 year |
| charlie | 50,000,000 | 5.0% | 50000000000000 | Vesting 4 years, cliff 1 year |
| **TOTAL** | **1,000,000,000** | **100%** | **1000000000000000** | |

**Summary:** Founders 20% · Community 80%.

For the full breakdown of distribution, vesting and multisig, see [docs/TOKENOMICS.md](docs/TOKENOMICS.md).

---

## Quick start

Requirements: Go 1.25+ and Ignite CLI v29.9.0 (Cosmos SDK v0.53.6).

### Build

```bash
git clone https://github.com/worrellchain/worrell.git
cd worrell
ignite chain build
```

This generates the `worrelld` binary.

### Init

```bash
worrelld init <moniker> --chain-id worrell-1
```

Node state is stored in `~/.worrell/`.

### Join testnet

```bash
# Initialize the node for the testnet
worrelld init <moniker> --chain-id worrell-testnet-1

# Download the testnet genesis, configure peers and start:
worrelld start
```

Full guide to operating a node and creating a validator: [docs/RUNNING-A-NODE.md](docs/RUNNING-A-NODE.md).

---

## Network parameters

| Parameter | Value |
|-----------|-------|
| Block time | ~5-6 seconds |
| Block gas limit | 40,000,000 |
| Min gas price | 0.025 uworrell |
| Pruning | Default |
| State Sync | Enabled |
| Telemetry | Prometheus enabled |

### Ports

| Service | Port |
|----------|-------:|
| CometBFT P2P | 26656 |
| CometBFT RPC | 26657 |
| REST API | 1317 |
| gRPC | 9090 |
| Prometheus | 26660 |
| Faucet (testnet only) | 4500 |

---

## Governance

Worrell supports **standard** and **expedited** (urgent) proposals. Both share the same quorum of **33.4%**.

| Parameter | Standard | Expedited |
|-----------|----------|-----------|
| Minimum deposit | 1,500,000,000 uworrell (1,500 WORRELL) | 7,500,000,000 uworrell (7,500 WORRELL) |
| Deposit period | 1,209,600 s (14 days) | 1,209,600 s (14 days) |
| Voting period | 432,000 s (5 days) | 86,400 s (24 hours) |
| Quorum | 33.4% (0.334) | 33.4% (0.334) |
| Threshold | 50% (0.50) | 66.7% (0.667) |
| Veto threshold | 33.4% (0.334) | 33.4% (0.334) |
| Burn on veto | true | true |
| Burn without quorum | true | true |

> Note: the Cosmos SDK `gov` module does not expose an `expedited_quorum` parameter; expedited proposals use the same quorum as standard ones.

Full detail: [docs/GOVERNANCE.md](docs/GOVERNANCE.md).

---

## Modules

| Module | Status | Configuration |
|--------|--------|---------------|
| Bank | Active | Transfers enabled |
| Staking | Active | See staking section |
| Governance | Active | Standard + expedited |
| Distribution | Active | Community tax 10% |
| Slashing | Active | See slashing section |
| Mint | Active | Dynamic inflation 7%-13% |
| Authz | Active | Standard module enabled |
| Fee Grants | Active | Standard module enabled |
| Upgrade | Active | Via governance proposals |
| IBC | Installed · disabled at genesis | `send_enabled: false`, `receive_enabled: false` |
| CosmWasm | Not installed | To be integrated post-launch via chain upgrade |
| Crisis | Not included | Deprecated in Cosmos SDK v0.53.6 |

---

## Staking

| Parameter | Mainnet | Testnet |
|-----------|---------|---------|
| Max validators | 100 | 100 |
| Min self-delegation | 1,000,000 uworrell (1 WORRELL) | 1,000,000 uworrell (1 WORRELL) |
| Minimum commission | 5% (0.05) | 5% (0.05) |
| Unbonding period | 1,814,400 s (21 days) | 3,600 s (1 hour) |
| Max entries | 7 | 7 |
| Historical entries | 10,000 | 10,000 |
| Bond denom | uworrell | uworrell |

> The 5% minimum commission is global. Each validator sets its own maximum commission.

---

## Slashing

| Event | Penalty | Additional effect (mainnet) |
|--------|--------------|----------------------------|
| Downtime | 0.01% (0.0001) | Jail 43,200 s (12 hours) |
| Double sign | 5% (0.05) | — |

| Parameter | Mainnet | Testnet |
|-----------|---------|---------|
| Signed blocks window | 10,000 blocks | 10,000 blocks |
| Min signed per window | 5% (0.05) | 5% (0.05) |
| Downtime jail duration | 43,200 s (12 hours) | 300 s (5 min) |

---

## Documentation

- [docs/TOKENOMICS.md](docs/TOKENOMICS.md) — Genesis distribution, vesting, multisig, inflation and fees.
- [docs/GOVERNANCE.md](docs/GOVERNANCE.md) — How governance works, proposal types and commands.
- [docs/RUNNING-A-NODE.md](docs/RUNNING-A-NODE.md) — Requirements, installation, operation and validator creation.

---

## Community

- **Validator announcements** (upgrades, governance, coordination): Telegram [t.me/worrellvalidators](https://t.me/worrellvalidators)
- **Questions and support:** [GitHub Discussions](https://github.com/worrellchain/worrell/discussions) · hello@worrellchain.com
- **News:** [@worrellchain](https://x.com/worrellchain) on X · [worrellchain.com](https://worrellchain.com)

---

## Security

If you find a vulnerability, please report it responsibly and privately to **security@worrellchain.com**. Please do not open public issues for security problems.

---

## Contributing

Contributions are welcome. Please review the [CONTRIBUTING.md](CONTRIBUTING.md) guide before opening a pull request.

---

## License

Distributed under the [Apache 2.0](https://github.com/worrellchain/worrell/blob/main/LICENSE) license.
