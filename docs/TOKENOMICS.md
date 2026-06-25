# Worrell Tokenomics

This document describes the economics of the **WORRELL** token: the total supply,
the genesis distribution, the founders' vesting schedule, the multisig
configuration of the custodied accounts, the treasury policy, inflation, and
transaction fees.

## Token and units

| Field | Value |
|-------|-------|
| Token | WORRELL |
| Base denomination | uworrell |
| Decimals | 6 |
| Conversion | 1 WORRELL = 1.000.000 uworrell |
| Total supply | 1.000.000.000 WORRELL (1.000.000.000.000.000 uworrell) |

---

## 1. Genesis distribution

The total supply is **1.000.000.000 WORRELL** and is split across 7 genesis
accounts. The founders receive 20 % combined and the community 80 %.

| Account | WORRELL | % | uworrell | Custody |
|--------|--------:|---:|---------:|----------|
| henry | 100.000.000 | 10 % | 100000000000000 | Vesting 4 years, cliff 1 year |
| george | 50.000.000 | 5 % | 50000000000000 | Vesting 4 years, cliff 1 year |
| charlie | 50.000.000 | 5 % | 50000000000000 | Vesting 4 years, cliff 1 year |
| treasury | 300.000.000 | 30 % | 300000000000000 | Multisig 3/3 |
| airdrop | 225.000.000 | 22,5 % | 225000000000000 | Multisig 2/4 |
| incentives | 175.000.000 | 17,5 % | 175000000000000 | Multisig 2/4 |
| reserve | 100.000.000 | 10 % | 100000000000000 | Multisig 2/4 |
| **TOTAL** | **1.000.000.000** | **100 %** | **1000000000000000** | |

**Summary:** Founders 20 % | Community 80 %

### Visual distribution

```
treasury    30.0%  ███████████████  300M
airdrop     22.5%  ███████████▎     225M
incentives  17.5%  ████████▊        175M
henry       10.0%  █████            100M
reserve     10.0%  █████            100M
george       5.0%  ██▌               50M
charlie      5.0%  ██▌               50M
                   └──────────────────────
                   Fundadores 20%  ·  Comunidad 80%
```

---

## 2. Founders and roles

The labels `henry`, `george`, and `charlie` identify genesis accounts and do not
represent specific individuals.

| Account | Allocation | Role |
|--------|-----------:|-----|
| henry | 100.000.000 WORRELL | Founder and initial network validator |
| george | 50.000.000 WORRELL | Founder and faucet source on testnet |
| charlie | 50.000.000 WORRELL | Founder |

All three founder accounts are subject to the same vesting scheme:
4 years in duration with a 1-year cliff.

---

## 3. Founder vesting

The allocations of the three founders are delivered through a
**`ContinuousVestingAccount`** type account.

> **Important:** the release is **linear and continuous**. There are no annual
> tranches or unlocks at discrete blocks. Once the cliff has passed, every
> instant of time releases a proportional fraction of the tokens.

### Vesting parameters

| Parameter | Value |
|-----------|-------|
| Account type | ContinuousVestingAccount |
| `start_time` | genesis timestamp + 31.557.600 s (1 year) |
| `end_time` | genesis timestamp + 126.230.400 s (4 years) |
| Year 0–1 (cliff) | 0 released: nothing transferable |
| Year 1–4 | Linear, continuous release over 3 years |

The `start_time` set one year after genesis implements the **cliff**: during the
first year absolutely nothing is released. From that moment onward, and up to the
`end_time` (4 years from genesis), the tokens are released in a linear and
continuous manner over the remaining 3 years.

### Release schedule — henry (100M)

The release is continuous; the points in the table are reference cuts to
illustrate the linear pace across the vesting span.

| Moment | Released in the interval | Cumulative |
|---------|-------------------------:|----------:|
| 0–12 months | 0 | 0 |
| 18 months | ~16.666.667 | ~16.666.667 |
| 24 months | ~16.666.667 | ~33.333.333 |
| 30 months | ~16.666.667 | ~50.000.000 |
| 36 months | ~16.666.667 | ~66.666.667 |
| 42 months | ~16.666.667 | ~83.333.333 |
| 48 months | ~16.666.667 | 100.000.000 |

### Release schedule — george and charlie (50M each)

| Moment | Cumulative |
|---------|----------:|
| 0–12 months | 0 |
| 24 months | ~16.666.667 |
| 36 months | ~33.333.333 |
| 48 months | 50.000.000 |

### Vesting rules

- Vesting tokens **CANNOT be transferred** until they are released.
- Vesting tokens **CAN be used for staking** (delegating).
- Vesting tokens **DO count for voting** in governance.
- Staking rewards are **fully LIQUID** and are not subject to vesting.

---

## 4. Multisig of custodied accounts

The four community accounts (treasury, airdrop, incentives, reserve) are
controlled by multisig addresses.

| Account | Type | Signers | Required signatures |
|--------|------|-----------|-------------------|
| treasury | 3 of 3 | henry, george, charlie | All 3 (unanimity) |
| airdrop | 2 of 4 | henry, george, charlie, airdrop-aux | 2 of the 3 founders |
| incentives | 2 of 4 | henry, george, charlie, incentives-aux | 2 of the 3 founders |
| reserve | 2 of 4 | henry, george, charlie, reserve-aux | 2 of the 3 founders |

> **Note on the auxiliary keys.** The keys `airdrop-aux`, `incentives-aux`, and
> `reserve-aux` exist **solely** to generate a multisig address distinct from
> treasury's. **They do not take part in the actual signing**: in practice, each
> of those accounts is operated with the signature of 2 of the 3 founders.

---

## 5. Treasury policy

Indicative (non-technical) distribution of the 300.000.000 WORRELL custodied in
the `treasury` account.

| Item | % of treasury | WORRELL |
|---------|--------------:|--------:|
| Development | 25 % | 75.000.000 |
| Liquidity | 20 % | 60.000.000 |
| Grants | 15 % | 45.000.000 |
| Partnerships | 15 % | 45.000.000 |
| Marketing | 15 % | 45.000.000 |
| Contingencies | 10 % | 30.000.000 |
| **TOTAL** | **100 %** | **300.000.000** |

---

## 6. Inflation and rewards

The mint module issues new tokens based on the percentage of the supply that is
staked, dynamically adjusting inflation.

| Parameter | Value |
|-----------|-------|
| Minimum inflation | 7 % |
| Maximum inflation | 13 % |
| Initial inflation | 13 % |
| Goal bonded ratio | 67 % |
| Mint denom | uworrell |

When less than 67 % of the supply is staked, inflation rises toward 13 %; when
more is staked, it falls toward 7 %.

### Reward distribution

| Destination | % of rewards |
|---------|---------------------:|
| Validators and delegators | 90 % |
| Community pool (community tax) | 10 % |

The community tax is **10 %**: that percentage of staking rewards is routed to
the community pool, and the remaining **90 %** is split between validators and
delegators.

---

## 7. Transaction fees

| Parameter | Value |
|-----------|-------|
| Min gas price | 0,025 uworrell |

Each node applies a minimum gas price of **0,025 uworrell**. Transactions that
offer a lower price will be rejected by the mempool.
