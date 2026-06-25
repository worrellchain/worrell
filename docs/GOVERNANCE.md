# Worrell Governance

Worrell is a proof-of-stake blockchain built with Cosmos SDK. Network
decision-making happens on-chain through the governance module (`x/gov`): any
token holder can propose changes, and the staker community approves or rejects
them by voting.

This document describes how a proposal's lifecycle works, the exact network
parameters, the deposit rules, the voting options, and the CLI commands needed
to participate.

- Token: **WORRELL**
- Base denomination: **uworrell** (1 WORRELL = 1,000,000 uworrell)
- Binary: **worrelld**
- Address prefix: **worrell1...**

---

## 1. How governance works

Every proposal goes through four sequential phases:

```
  submit  ──▶  deposit  ──▶  vote  ──▶  outcome
```

1. **Submit.** A user submits the proposal to the chain. At submission time they
   can attach an initial deposit (partial or full). The proposal is recorded
   with an identifier (`proposal-id`).

2. **Deposit (deposit period).** The proposal must reach the **minimum deposit**
   before the deposit period ends. Any account can contribute to a proposal's
   deposit, not just its creator.
   - If the minimum deposit is reached within the time limit, the proposal
     automatically moves to voting.
   - If the minimum deposit is **not** reached within the time limit, the
     proposal expires and the deposits are returned (not burned, because
     `burn_proposal_deposit_prevote` is `false`).

3. **Vote (voting period).** Begins once the minimum deposit is covered. During
   this period, stakers vote with one of four options: `Yes`, `No`,
   `NoWithVeto`, `Abstain`. The weight of each vote is proportional to staked
   tokens (vesting tokens also count for voting).

4. **Outcome.** When voting ends, the result is computed from the **quorum**,
   the **threshold**, and the **veto threshold**:
   - For the result to be valid, participation must reach the **quorum**.
   - If the quorum is reached and `Yes` votes exceed the **threshold** (without
     exceeding the veto threshold), the proposal is **approved** and, where
     applicable, executed automatically.
   - In any other case, the proposal is **rejected**.

The specific rules for returning or burning the deposit depending on the outcome
are detailed in section [4. Deposit rules](#4-deposit-rules).

---

## 2. Standard proposals

Governance module parameters for standard proposals. The **Testnet** column
shows the reduced values used in `worrell-testnet-1` so the full cycle can be
tested in minutes instead of days.

| Parameter | Mainnet value | Testnet value |
|-----------|---------------|---------------|
| Minimum deposit | **1,500,000,000 uworrell** (1,500 WORRELL) | Same (1,500,000,000 uworrell) |
| Deposit period | 1,209,600s (14 days) | 600s (10 min) |
| Voting period | 432,000s (5 days) | 300s (5 min) |
| Quorum | 33.4% (0.334) | Same |
| Threshold | 50% (0.50) | Same |
| Veto threshold | 33.4% (0.334) | Same |
| Burn on veto | `true` | `true` |
| Burn without quorum | `true` | `true` |
| Burn before voting | `false` (returned) | `false` |

---

## 3. Expedited proposals (urgent)

**Expedited** proposals allow urgent matters to be resolved with a shorter
voting period, in exchange for a higher minimum deposit and a more demanding
approval threshold.

| Parameter | Mainnet value | Testnet value |
|-----------|---------------|---------------|
| Minimum deposit | **7,500,000,000 uworrell** (7,500 WORRELL) | Same (7,500,000,000 uworrell) |
| Voting period | 86,400s (24 hours) | 120s (2 min) |
| Threshold | 66.7% (0.667) | Same |
| Quorum | **33.4% (0.334) — the same as standard ones** | Same |

> ### ⚠️ IMPORTANT: there is no `expedited_quorum`
>
> Expedited proposals use **exactly the same quorum as standard ones:
> 33.4% (0.334)**.
>
> The Cosmos SDK `x/gov` module does **NOT** have an `expedited_quorum`
> parameter. What the module allows to be configured independently for expedited
> proposals is only the **minimum deposit**, the **voting period**, and the
> **threshold**. The quorum is a global parameter shared by both proposal types.
>
> Do not attempt to configure a different quorum for expedited proposals: that
> parameter does not exist. Any proposal (standard or expedited) that does not
> reach 33.4% participation will not pass the quorum.

---

## 4. Voting options

Each voter chooses one of these four options:

| Option | Meaning | Counts toward quorum |
|--------|---------|----------------------|
| `Yes` | In favor of the proposal | Yes |
| `No` | Against the proposal | Yes |
| `NoWithVeto` | Against **and** considers the proposal abusive or spam; counts toward the veto threshold | Yes |
| `Abstain` | No position; adds to participation but neither for nor against | Yes |

Notes:

- `Abstain` **counts toward the quorum** (participation) even though it expresses
  no position for or against.
- `NoWithVeto` is an especially strong vote: if the sum of `NoWithVeto` exceeds
  the **veto threshold (33.4%)** of the votes cast, the proposal is rejected by
  veto **and its deposit is burned** (see section 5), even if the `Yes` votes
  had exceeded the threshold.

---

## 5. Deposit rules

The fate of the deposit depends on the proposal's outcome. Worrell uses this
configuration:

| Parameter | Value | Effect |
|-----------|-------|--------|
| `burn_vote_quorum` | `true` | The deposit **is burned** if the vote **does not reach the quorum** (33.4%) |
| `burn_proposal_deposit_prevote` | `false` | The deposit is **NOT burned** if the proposal does not reach voting (deposit period expired): it is **returned** |
| Burn on veto (`NoWithVeto` exceeds the threshold) | `true` | The deposit **is burned** |

Summary of the deposit's fate by scenario:

| Scenario | Deposit's fate |
|----------|----------------|
| Proposal **approved** | **Returned** to depositors |
| Proposal **rejected normally** (`No` wins, no veto) | **Returned** to depositors |
| Proposal **vetoed** (`NoWithVeto` exceeds the 33.4% threshold) | **BURNED** |
| Vote **without quorum** (`burn_vote_quorum: true`) | **BURNED** |
| Minimum deposit **not reached** within the deposit period (`burn_proposal_deposit_prevote: false`) | **Returned** (not burned before voting) |

Key points:

- The deposit is **returned** in a normal rejection (the proposal lost, but was
  not vetoed and there was a quorum).
- The deposit is **burned** in two cases: (1) when there is a **veto** because
  `NoWithVeto` exceeds its threshold, and (2) when the **quorum is not reached**,
  because `burn_vote_quorum` is set to `true`.
- The deposit is **not burned before voting**: if the proposal never reaches a
  vote due to lack of minimum deposit, the funds are returned
  (`burn_proposal_deposit_prevote: false`).

---

## 6. Proposal types

Worrell's governance module supports the following proposal types:

| Type | Description |
|------|-------------|
| **Parameter change** | Modifies on-chain parameters of a module (staking, slashing, mint, distribution, gov, etc.). |
| **Software upgrade** | Schedules a coordinated upgrade of the `worrelld` binary at a specific block height (`upgrade` module). |
| **Community pool spend** | Spends funds from the community pool (fed by the 10% community tax) to a destination address. |
| **Text** | Signaling proposal with no automatic on-chain effect; used to gauge community sentiment. |

> In Cosmos SDK v0.50+ / v0.53, these types are usually expressed as messages
> (`MsgUpdateParams`, `MsgSoftwareUpgrade`, `MsgCommunityPoolSpend`, etc.) within
> a generic proposal (`submit-proposal` with a JSON file). A **text** proposal is
> a proposal with no messages to execute it.

---

## 7. Example proposal in JSON

Example `proposal.json` file for a parameter change proposal. It includes an
**initial deposit of `1500000000uworrell`** (the standard minimum deposit of
1,500 WORRELL):

```json
{
  "messages": [
    {
      "@type": "/cosmos.staking.v1beta1.MsgUpdateParams",
      "authority": "worrell10d07y265gmmuvt4z0w9aw880jnsr700jp5y0fr",
      "params": {
        "unbonding_time": "1814400s",
        "max_validators": 100,
        "max_entries": 7,
        "historical_entries": 10000,
        "bond_denom": "uworrell",
        "min_commission_rate": "0.050000000000000000"
      }
    }
  ],
  "metadata": "ipfs://CID",
  "deposit": "1500000000uworrell",
  "title": "Staking parameter adjustment",
  "summary": "Example proposal that updates the staking module parameters.",
  "expedited": false
}
```

Notes about the JSON:

- The `deposit` field uses the base denomination: `"1500000000uworrell"`
  (= 1,500 WORRELL), which covers the standard minimum deposit.
- The `authority` field is the governance module account (the `gov` module
  address); the proposal can only be executed through governance.
- For an **expedited** proposal, change `"expedited": true` and remember that the
  minimum deposit is `7500000000uworrell` (7,500 WORRELL).
- A **text** proposal carries `"messages": []` (no messages to execute).

---

## 8. CLI commands

All commands use the **`worrelld`** binary, the **`uworrell`** base denomination,
and addresses with the **`worrell1...`** prefix.

### 8.1 Submit a proposal

From a `proposal.json` file like the one in the previous section:

```bash
worrelld tx gov submit-proposal proposal.json \
  --from worrell1tuxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  --chain-id worrell-1 \
  --gas auto \
  --gas-adjustment 1.5 \
  --gas-prices 0.025uworrell \
  --yes
```

To submit an additional deposit to an existing proposal during the deposit
period:

```bash
worrelld tx gov deposit 1 1500000000uworrell \
  --from worrell1tuxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  --chain-id worrell-1 \
  --gas-prices 0.025uworrell \
  --yes
```

> `1` is the `proposal-id`. On testnet use `--chain-id worrell-testnet-1`.

### 8.2 Query proposals

```bash
# List all proposals
worrelld query gov proposals

# View a specific proposal (proposal-id = 1)
worrelld query gov proposal 1

# View the vote tally
worrelld query gov tally 1

# View a proposal's deposits
worrelld query gov deposits 1

# View a specific account's vote
worrelld query gov vote 1 worrell1tuxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# View the current governance parameters
worrelld query gov params
```

### 8.3 Vote

During the voting period, with one of the options `yes`, `no`, `no_with_veto`,
or `abstain`:

```bash
worrelld tx gov vote 1 yes \
  --from worrell1tuxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  --chain-id worrell-1 \
  --gas-prices 0.025uworrell \
  --yes
```

Examples of the four voting options:

```bash
worrelld tx gov vote 1 yes          --from <account> --chain-id worrell-1 --gas-prices 0.025uworrell --yes
worrelld tx gov vote 1 no           --from <account> --chain-id worrell-1 --gas-prices 0.025uworrell --yes
worrelld tx gov vote 1 no_with_veto --from <account> --chain-id worrell-1 --gas-prices 0.025uworrell --yes
worrelld tx gov vote 1 abstain      --from <account> --chain-id worrell-1 --gas-prices 0.025uworrell --yes
```

Weighted vote, to split the weight across several options:

```bash
worrelld tx gov weighted-vote 1 yes=0.7,abstain=0.3 \
  --from worrell1tuxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  --chain-id worrell-1 \
  --gas-prices 0.025uworrell \
  --yes
```

---

## 9. Quick summary

- **Cycle:** submit → deposit → vote → outcome.
- **Standard minimum deposit:** 1,500 WORRELL (`1500000000uworrell`).
  **Expedited:** 7,500 WORRELL (`7500000000uworrell`).
- **Quorum:** 33.4% for **both** types. **There is no `expedited_quorum`.**
- **Threshold:** 50% standard, 66.7% expedited. **Veto:** 33.4%.
- **Deposit:** returned on approval and on normal rejection; **burned** on veto
  and on lack of quorum; not burned before voting.
- **Voting options:** `Yes`, `No`, `NoWithVeto`, `Abstain`.
- **Types:** parameter change, software upgrade, community pool spend, text.
