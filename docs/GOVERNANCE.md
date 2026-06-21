# Gobernanza de Worrell

Worrell es una blockchain proof-of-stake construida con Cosmos SDK. La toma de
decisiones de la red se realiza on-chain mediante el módulo de gobernanza
(`x/gov`): cualquier titular de tokens puede proponer cambios y la comunidad de
stakers los aprueba o rechaza votando.

Este documento describe cómo funciona el ciclo de una propuesta, los parámetros
exactos de la red, las reglas de depósito, las opciones de voto y los comandos
de CLI necesarios para participar.

- Token: **WORRELL**
- Denominación base: **uworrell** (1 WORRELL = 1,000,000 uworrell)
- Binario: **worrelld**
- Prefijo de direcciones: **worrell1...**

---

## 1. Cómo funciona la gobernanza

Toda propuesta atraviesa cuatro fases secuenciales:

```
  submit  ──▶  deposit  ──▶  vote  ──▶  outcome
 (enviar)     (depósito)    (votar)    (resultado)
```

1. **Submit (envío).** Un usuario envía la propuesta a la cadena. En el momento
   del envío puede adjuntar un depósito inicial (parcial o total). La propuesta
   queda registrada con un identificador (`proposal-id`).

2. **Deposit (periodo de depósito).** La propuesta debe alcanzar el **depósito
   mínimo** antes de que termine el periodo de depósito. Cualquier cuenta puede
   contribuir al depósito de una propuesta, no solo quien la creó.
   - Si se alcanza el depósito mínimo dentro del plazo, la propuesta pasa
     automáticamente a votación.
   - Si **no** se alcanza el depósito mínimo dentro del plazo, la propuesta
     caduca y los depósitos se devuelven (no se queman, porque
     `burn_proposal_deposit_prevote` es `false`).

3. **Vote (periodo de votación).** Comienza cuando se cubre el depósito mínimo.
   Durante este periodo, los stakers votan con una de cuatro opciones: `Yes`,
   `No`, `NoWithVeto`, `Abstain`. El peso de cada voto es proporcional a los
   tokens en staking (los tokens en vesting también cuentan para votar).

4. **Outcome (resultado).** Al terminar la votación, el resultado se calcula a
   partir del **quórum**, el **threshold** y el **veto threshold**:
   - Para que el resultado sea válido, la participación debe alcanzar el
     **quórum**.
   - Si se alcanza el quórum y los votos `Yes` superan el **threshold** (sin
     superar el umbral de veto), la propuesta se **aprueba** y, en su caso, se
     ejecuta automáticamente.
   - En cualquier otro caso, la propuesta se **rechaza**.

Las reglas concretas de devolución o quema del depósito según el resultado se
detallan en la sección [4. Reglas de depósito](#4-reglas-de-depósito).

---

## 2. Propuestas estándar

Parámetros del módulo de gobernanza para propuestas estándar. La columna
**Testnet** indica los valores reducidos usados en `worrell-testnet-1` para
poder probar el ciclo completo en minutos en lugar de días.

| Parámetro | Valor mainnet | Valor testnet |
|-----------|---------------|---------------|
| Depósito mínimo | **1,500,000,000 uworrell** (1,500 WORRELL) | Igual (1,500,000,000 uworrell) |
| Periodo de depósito | 1,209,600s (14 días) | 600s (10 min) |
| Periodo de votación | 432,000s (5 días) | 300s (5 min) |
| Quórum | 33.4% (0.334) | Igual |
| Threshold | 50% (0.50) | Igual |
| Veto threshold | 33.4% (0.334) | Igual |
| Burn en veto | `true` | `true` |
| Burn sin quórum | `true` | `true` |
| Burn antes de votación | `false` (se devuelve) | `false` |

---

## 3. Propuestas expedited (urgentes)

Las propuestas **expedited** permiten resolver asuntos urgentes con un periodo
de votación más corto, a cambio de un depósito mínimo mayor y un threshold de
aprobación más exigente.

| Parámetro | Valor mainnet | Valor testnet |
|-----------|---------------|---------------|
| Depósito mínimo | **7,500,000,000 uworrell** (7,500 WORRELL) | Igual (7,500,000,000 uworrell) |
| Periodo de votación | 86,400s (24 horas) | 120s (2 min) |
| Threshold | 66.7% (0.667) | Igual |
| Quórum | **33.4% (0.334) — el mismo que las estándar** | Igual |

> ### ⚠️ IMPORTANTE: no existe `expedited_quorum`
>
> Las propuestas expedited usan **exactamente el mismo quórum que las
> estándar: 33.4% (0.334)**.
>
> El módulo `x/gov` de Cosmos SDK **NO** tiene un parámetro
> `expedited_quorum`. Lo que el módulo permite ajustar de forma independiente
> para las expedited es únicamente el **depósito mínimo**, el **periodo de
> votación** y el **threshold**. El quórum es un parámetro global compartido por
> ambos tipos de propuesta.
>
> No intentes configurar un quórum distinto para las expedited: ese parámetro no
> existe. Cualquier propuesta (estándar o expedited) que no alcance el 33.4% de
> participación no superará el quórum.

---

## 4. Opciones de voto

Cada votante elige una de estas cuatro opciones:

| Opción | Significado | Cuenta para quórum |
|--------|-------------|--------------------|
| `Yes` | A favor de la propuesta | Sí |
| `No` | En contra de la propuesta | Sí |
| `NoWithVeto` | En contra **y** considera la propuesta abusiva o spam; cuenta para el umbral de veto | Sí |
| `Abstain` | Sin posición; suma a la participación pero no a favor ni en contra | Sí |

Notas:

- `Abstain` **cuenta para el quórum** (participación) aunque no exprese una
  posición a favor o en contra.
- `NoWithVeto` es un voto especialmente fuerte: si la suma de `NoWithVeto`
  supera el **veto threshold (33.4%)** de los votos emitidos, la propuesta se
  rechaza por veto **y su depósito se quema** (ver sección 5), aunque los `Yes`
  hubieran superado el threshold.

---

## 5. Reglas de depósito

El destino del depósito depende del resultado de la propuesta. Worrell usa esta
configuración:

| Parámetro | Valor | Efecto |
|-----------|-------|--------|
| `burn_vote_quorum` | `true` | **Se quema** el depósito si la votación **no alcanza el quórum** (33.4%) |
| `burn_proposal_deposit_prevote` | `false` | **NO se quema** el depósito si la propuesta no llega a votación (deposit period expirado): se **devuelve** |
| Burn en veto (`NoWithVeto` supera el umbral) | `true` | **Se quema** el depósito |

Resumen del destino del depósito según el escenario:

| Escenario | Destino del depósito |
|-----------|----------------------|
| Propuesta **aprobada** | Se **devuelve** a los depositantes |
| Propuesta **rechazada de forma normal** (gana `No`, sin veto) | Se **devuelve** a los depositantes |
| Propuesta **vetada** (`NoWithVeto` supera el umbral del 33.4%) | Se **QUEMA** |
| Votación **sin quórum** (`burn_vote_quorum: true`) | Se **QUEMA** |
| Depósito mínimo **no alcanzado** dentro del periodo de depósito (`burn_proposal_deposit_prevote: false`) | Se **devuelve** (no se quema antes de votación) |

Puntos clave:

- El depósito se **devuelve** en un rechazo normal (la propuesta perdió, pero no
  fue vetada y sí hubo quórum).
- El depósito se **quema** en dos casos: (1) cuando hay **veto** porque
  `NoWithVeto` supera su umbral, y (2) cuando **no se alcanza el quórum**, porque
  `burn_vote_quorum` está en `true`.
- El depósito **no se quema antes de la votación**: si la propuesta nunca llega a
  votarse por falta de depósito mínimo, los fondos se devuelven
  (`burn_proposal_deposit_prevote: false`).

---

## 6. Tipos de propuestas

El módulo de gobernanza de Worrell admite los siguientes tipos de propuesta:

| Tipo | Descripción |
|------|-------------|
| **Parameter change** | Modifica parámetros on-chain de un módulo (staking, slashing, mint, distribution, gov, etc.). |
| **Software upgrade** | Programa una actualización coordinada del binario `worrelld` a una altura de bloque concreta (módulo `upgrade`). |
| **Community pool spend** | Gasta fondos del community pool (alimentado por el 10% de community tax) hacia una dirección destino. |
| **Text** | Propuesta de señalización sin efecto on-chain automático; sirve para medir el sentir de la comunidad. |

> En Cosmos SDK v0.50+ / v0.53, estos tipos se expresan habitualmente como
> mensajes (`MsgUpdateParams`, `MsgSoftwareUpgrade`,
> `MsgCommunityPoolSpend`, etc.) dentro de una propuesta genérica
> (`submit-proposal` con un archivo JSON). Una propuesta de tipo **text** es una
> propuesta sin mensajes que la ejecuten.

---

## 7. Ejemplo de propuesta en JSON

Ejemplo de archivo `proposal.json` para una propuesta de cambio de parámetro.
Incluye un **depósito inicial de `1500000000uworrell`** (el depósito mínimo
estándar de 1,500 WORRELL):

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
  "title": "Ajuste de parámetros de staking",
  "summary": "Propuesta de ejemplo que actualiza los parámetros del módulo de staking.",
  "expedited": false
}
```

Notas sobre el JSON:

- El campo `deposit` usa la denominación base: `"1500000000uworrell"`
  (= 1,500 WORRELL), que cubre el depósito mínimo estándar.
- El campo `authority` es la cuenta del módulo de gobernanza (la dirección del
  módulo `gov`); la propuesta solo puede ejecutarse a través de gobernanza.
- Para una propuesta **expedited**, cambia `"expedited": true` y recuerda que el
  depósito mínimo es de `7500000000uworrell` (7,500 WORRELL).
- Una propuesta de tipo **text** lleva `"messages": []` (sin mensajes a
  ejecutar).

---

## 8. Comandos de CLI

Todos los comandos usan el binario **`worrelld`**, la denominación base
**`uworrell`** y direcciones con prefijo **`worrell1...`**.

### 8.1 Enviar una propuesta (submit)

A partir de un archivo `proposal.json` como el de la sección anterior:

```bash
worrelld tx gov submit-proposal proposal.json \
  --from worrell1tuxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  --chain-id worrell-1 \
  --gas auto \
  --gas-adjustment 1.5 \
  --gas-prices 0.025uworrell \
  --yes
```

Para enviar un depósito adicional a una propuesta existente durante el periodo
de depósito:

```bash
worrelld tx gov deposit 1 1500000000uworrell \
  --from worrell1tuxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  --chain-id worrell-1 \
  --gas-prices 0.025uworrell \
  --yes
```

> `1` es el `proposal-id`. En testnet usa `--chain-id worrell-testnet-1`.

### 8.2 Consultar propuestas (query)

```bash
# Listar todas las propuestas
worrelld query gov proposals

# Ver una propuesta concreta (proposal-id = 1)
worrelld query gov proposal 1

# Ver el resultado de la votación
worrelld query gov tally 1

# Ver los depósitos de una propuesta
worrelld query gov deposits 1

# Ver el voto de una cuenta concreta
worrelld query gov vote 1 worrell1tuxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# Ver los parámetros de gobernanza vigentes
worrelld query gov params
```

### 8.3 Votar (vote)

Durante el periodo de votación, con una de las opciones `yes`, `no`,
`no_with_veto` o `abstain`:

```bash
worrelld tx gov vote 1 yes \
  --from worrell1tuxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  --chain-id worrell-1 \
  --gas-prices 0.025uworrell \
  --yes
```

Ejemplos de las cuatro opciones de voto:

```bash
worrelld tx gov vote 1 yes          --from <cuenta> --chain-id worrell-1 --gas-prices 0.025uworrell --yes
worrelld tx gov vote 1 no           --from <cuenta> --chain-id worrell-1 --gas-prices 0.025uworrell --yes
worrelld tx gov vote 1 no_with_veto --from <cuenta> --chain-id worrell-1 --gas-prices 0.025uworrell --yes
worrelld tx gov vote 1 abstain      --from <cuenta> --chain-id worrell-1 --gas-prices 0.025uworrell --yes
```

Voto ponderado (weighted vote), para repartir el peso entre varias opciones:

```bash
worrelld tx gov weighted-vote 1 yes=0.7,abstain=0.3 \
  --from worrell1tuxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  --chain-id worrell-1 \
  --gas-prices 0.025uworrell \
  --yes
```

---

## 9. Resumen rápido

- **Ciclo:** submit → deposit → vote → outcome.
- **Depósito mínimo estándar:** 1,500 WORRELL (`1500000000uworrell`).
  **Expedited:** 7,500 WORRELL (`7500000000uworrell`).
- **Quórum:** 33.4% para **ambos** tipos. **No existe `expedited_quorum`.**
- **Threshold:** 50% estándar, 66.7% expedited. **Veto:** 33.4%.
- **Depósito:** se devuelve en aprobación y en rechazo normal; se **quema** con
  veto y por falta de quórum; no se quema antes de votación.
- **Opciones de voto:** `Yes`, `No`, `NoWithVeto`, `Abstain`.
- **Tipos:** parameter change, software upgrade, community pool spend, text.
