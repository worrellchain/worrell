# Worrell

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://github.com/worrellchain/worrell/blob/main/LICENSE)

**Worrell** es una blockchain proof-of-stake (PoS) construida con [Ignite CLI](https://ignite.com/) y [Cosmos SDK](https://docs.cosmos.network/), enfocada en **pagos e infraestructura energética**.

| Campo | Valor |
|-------|-------|
| Token | WORRELL |
| Denominación base | uworrell |
| Decimales | 6 (1 WORRELL = 1.000.000 uworrell) |
| Binario | `worrelld` |
| Data directory | `~/.worrell/` |
| Address prefix | `worrell` (direcciones `worrell1...`) |
| Chain ID mainnet | `worrell-1` |
| Chain ID testnet | `worrell-testnet-1` |
| Licencia | Apache 2.0 |

---

## Características principales

- **Consenso Proof-of-Stake** sobre CometBFT, con tiempos de bloque de ~5-6 segundos.
- **Enfoque en pagos e infraestructura energética** como casos de uso prioritarios.
- **Gobernanza on-chain** con propuestas estándar y expedited (urgentes).
- **Staking y delegación** con hasta 100 validadores activos.
- **Inflación dinámica** ligada al ratio de tokens en staking (objetivo 67%).
- **Vesting de fundadores** lineal y continuo, con cuentas multisig para la tesorería y las reservas de la comunidad.
- **IBC instalado** (desactivado en genesis, se habilitará por gobernanza cuando la red sea estable).

---

## Token economics

**Supply total: 1.000.000.000 WORRELL** (1.000.000.000.000.000 uworrell).

Distribución genesis:

```
treasury    30.0%  ██████████████████████████████
airdrop     22.5%  ██████████████████████▌
incentives  17.5%  █████████████████▌
henry       10.0%  ██████████
reserve     10.0%  ██████████
george       5.0%  █████
charlie      5.0%  █████
```

| Cuenta | WORRELL | % | uworrell | Custodia |
|--------|--------:|--:|---------:|----------|
| treasury | 300.000.000 | 30,0% | 300000000000000 | Multisig 3/3 |
| airdrop | 225.000.000 | 22,5% | 225000000000000 | Multisig 2/4 |
| incentives | 175.000.000 | 17,5% | 175000000000000 | Multisig 2/4 |
| henry | 100.000.000 | 10,0% | 100000000000000 | Vesting 4 años, cliff 1 año |
| reserve | 100.000.000 | 10,0% | 100000000000000 | Multisig 2/4 |
| george | 50.000.000 | 5,0% | 50000000000000 | Vesting 4 años, cliff 1 año |
| charlie | 50.000.000 | 5,0% | 50000000000000 | Vesting 4 años, cliff 1 año |
| **TOTAL** | **1.000.000.000** | **100%** | **1000000000000000** | |

**Resumen:** Fundadores 20% · Comunidad 80%.

Para el detalle completo de distribución, vesting y multisig consulta [docs/TOKENOMICS.md](docs/TOKENOMICS.md).

---

## Quick start

Requisitos: Go 1.22+ e Ignite CLI v29.9.0 (Cosmos SDK v0.53.6).

### Build

```bash
git clone https://github.com/worrellchain/worrell.git
cd worrell
ignite chain build
```

Esto genera el binario `worrelld`.

### Init

```bash
worrelld init <moniker> --chain-id worrell-1
```

El estado del nodo se almacena en `~/.worrell/`.

### Join testnet

```bash
# Inicializa el nodo para la testnet
worrelld init <moniker> --chain-id worrell-testnet-1

# Descarga el genesis de la testnet, configura los peers y arranca:
worrelld start
```

Guía completa para operar un nodo y crear un validador: [docs/RUNNING-A-NODE.md](docs/RUNNING-A-NODE.md).

---

## Parámetros de red

| Parámetro | Valor |
|-----------|-------|
| Block time | ~5-6 segundos |
| Block gas limit | 40.000.000 |
| Min gas price | 0.025 uworrell |
| Pruning | Default |
| State Sync | Habilitado |
| Telemetría | Prometheus habilitado |

### Puertos

| Servicio | Puerto |
|----------|-------:|
| CometBFT P2P | 26656 |
| CometBFT RPC | 26657 |
| API REST | 1317 |
| gRPC | 9090 |
| Prometheus | 26660 |
| Faucet (solo testnet) | 4500 |

---

## Gobernanza

Worrell soporta propuestas **estándar** y **expedited** (urgentes). Ambas comparten el mismo quórum del **33,4%**.

| Parámetro | Estándar | Expedited |
|-----------|----------|-----------|
| Depósito mínimo | 1.500.000.000 uworrell (1.500 WORRELL) | 7.500.000.000 uworrell (7.500 WORRELL) |
| Periodo de depósito | 1.209.600 s (14 días) | 1.209.600 s (14 días) |
| Periodo de votación | 432.000 s (5 días) | 86.400 s (24 horas) |
| Quórum | 33,4% (0.334) | 33,4% (0.334) |
| Threshold | 50% (0.50) | 66,7% (0.667) |
| Veto threshold | 33,4% (0.334) | 33,4% (0.334) |
| Burn en veto | true | true |
| Burn sin quórum | true | true |

> Nota: el módulo `gov` de Cosmos SDK no expone un parámetro `expedited_quorum`; las propuestas expedited usan el mismo quórum que las estándar.

Detalle completo: [docs/GOVERNANCE.md](docs/GOVERNANCE.md).

---

## Módulos

| Módulo | Estado | Configuración |
|--------|--------|---------------|
| Bank | Activo | Transferencias habilitadas |
| Staking | Activo | Ver sección de staking |
| Governance | Activo | Estándar + expedited |
| Distribution | Activo | Community tax 10% |
| Slashing | Activo | Ver sección de slashing |
| Mint | Activo | Inflación dinámica 7%-13% |
| Authz | Activo | Módulo estándar habilitado |
| Fee Grants | Activo | Módulo estándar habilitado |
| Upgrade | Activo | Vía propuestas de gobernanza |
| IBC | Instalado · desactivado en genesis | `send_enabled: false`, `receive_enabled: false` |
| CosmWasm | No instalado | Se integrará post-lanzamiento vía chain upgrade |
| Crisis | No incluido | Deprecado en Cosmos SDK v0.53.6 |

---

## Staking

| Parámetro | Mainnet | Testnet |
|-----------|---------|---------|
| Max validadores | 100 | 100 |
| Min self-delegation | 1.000.000 uworrell (1 WORRELL) | 1.000.000 uworrell (1 WORRELL) |
| Comisión mínima | 5% (0.05) | 5% (0.05) |
| Unbonding period | 1.814.400 s (21 días) | 3.600 s (1 hora) |
| Max entries | 7 | 7 |
| Historical entries | 10.000 | 10.000 |
| Bond denom | uworrell | uworrell |

> La comisión mínima del 5% es global. Cada validador define su propio máximo de comisión.

---

## Slashing

| Evento | Penalización | Efecto adicional (mainnet) |
|--------|--------------|----------------------------|
| Downtime | 0,01% (0.0001) | Jail 43.200 s (12 horas) |
| Double sign | 5% (0.05) | — |

| Parámetro | Mainnet | Testnet |
|-----------|---------|---------|
| Signed blocks window | 10.000 bloques | 10.000 bloques |
| Min signed per window | 5% (0.05) | 5% (0.05) |
| Downtime jail duration | 43.200 s (12 horas) | 300 s (5 min) |

---

## Documentación

- [docs/TOKENOMICS.md](docs/TOKENOMICS.md) — Distribución genesis, vesting, multisig, inflación y fees.
- [docs/GOVERNANCE.md](docs/GOVERNANCE.md) — Cómo funciona la gobernanza, tipos de propuesta y comandos.
- [docs/RUNNING-A-NODE.md](docs/RUNNING-A-NODE.md) — Requisitos, instalación, operación y creación de validador.

---

## Seguridad

Si encuentras una vulnerabilidad, repórtala de forma responsable y privada a **security@worrellchain.io**. Por favor, no abras issues públicas para problemas de seguridad.

---

## Contribuir

Las contribuciones son bienvenidas. Revisa la guía [CONTRIBUTING.md](CONTRIBUTING.md) antes de abrir un pull request.

---

## Licencia

Distribuido bajo la licencia [Apache 2.0](https://github.com/worrellchain/worrell/blob/main/LICENSE).
