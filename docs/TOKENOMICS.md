# Tokenomics de Worrell

Este documento describe la economía del token **WORRELL**: el supply total, la
distribución de génesis, el calendario de vesting de los fundadores, la
configuración multisig de las cuentas custodiadas, la política de tesorería, la
inflación y los fees de transacción.

## Token y unidades

| Campo | Valor |
|-------|-------|
| Token | WORRELL |
| Denominación base | uworrell |
| Decimales | 6 |
| Conversión | 1 WORRELL = 1.000.000 uworrell |
| Supply total | 1.000.000.000 WORRELL (1.000.000.000.000.000 uworrell) |

---

## 1. Distribución de génesis

El supply total es de **1.000.000.000 WORRELL** y se reparte en 7 cuentas de
génesis. Los fundadores reciben en conjunto el 20 % y la comunidad el 80 %.

| Cuenta | WORRELL | % | uworrell | Custodia |
|--------|--------:|---:|---------:|----------|
| henry | 100.000.000 | 10 % | 100000000000000 | Vesting 4 años, cliff 1 año |
| george | 50.000.000 | 5 % | 50000000000000 | Vesting 4 años, cliff 1 año |
| charlie | 50.000.000 | 5 % | 50000000000000 | Vesting 4 años, cliff 1 año |
| treasury | 300.000.000 | 30 % | 300000000000000 | Multisig 3/3 |
| airdrop | 225.000.000 | 22,5 % | 225000000000000 | Multisig 2/4 |
| incentives | 175.000.000 | 17,5 % | 175000000000000 | Multisig 2/4 |
| reserve | 100.000.000 | 10 % | 100000000000000 | Multisig 2/4 |
| **TOTAL** | **1.000.000.000** | **100 %** | **1000000000000000** | |

**Resumen:** Fundadores 20 % | Comunidad 80 %

### Distribución visual

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

## 2. Fundadores y roles

Las etiquetas `henry`, `george` y `charlie` identifican cuentas de génesis y no
representan personas concretas.

| Cuenta | Asignación | Rol |
|--------|-----------:|-----|
| henry | 100.000.000 WORRELL | Fundador y validador inicial de la red |
| george | 50.000.000 WORRELL | Fundador y fuente del faucet en testnet |
| charlie | 50.000.000 WORRELL | Fundador |

Las tres cuentas de fundadores están sujetas al mismo esquema de vesting:
4 años de duración con un cliff de 1 año.

---

## 3. Vesting de fundadores

Las asignaciones de los tres fundadores se entregan mediante una cuenta de tipo
**`ContinuousVestingAccount`**.

> **Importante:** la liberación es **lineal y continua**. No existen tramos
> anuales ni desbloqueos por bloques discretos. Una vez superado el cliff, cada
> instante de tiempo libera una fracción proporcional de los tokens.

### Parámetros del vesting

| Parámetro | Valor |
|-----------|-------|
| Tipo de cuenta | ContinuousVestingAccount |
| `start_time` | timestamp de génesis + 31.557.600 s (1 año) |
| `end_time` | timestamp de génesis + 126.230.400 s (4 años) |
| Año 0–1 (cliff) | 0 liberado: nada transferible |
| Año 1–4 | Liberación lineal continua durante 3 años |

El `start_time` situado un año después del génesis implementa el **cliff**:
durante el primer año no se libera absolutamente nada. A partir de ese momento, y
hasta el `end_time` (4 años desde el génesis), los tokens se liberan de forma
lineal y continua a lo largo de los 3 años restantes.

### Calendario de liberación — henry (100M)

La liberación es continua; los puntos de la tabla son cortes de referencia para
ilustrar el ritmo lineal a lo largo del tramo de vesting.

| Momento | Liberado en el intervalo | Acumulado |
|---------|-------------------------:|----------:|
| 0–12 meses | 0 | 0 |
| 18 meses | ~16.666.667 | ~16.666.667 |
| 24 meses | ~16.666.667 | ~33.333.333 |
| 30 meses | ~16.666.667 | ~50.000.000 |
| 36 meses | ~16.666.667 | ~66.666.667 |
| 42 meses | ~16.666.667 | ~83.333.333 |
| 48 meses | ~16.666.667 | 100.000.000 |

### Calendario de liberación — george y charlie (50M cada uno)

| Momento | Acumulado |
|---------|----------:|
| 0–12 meses | 0 |
| 24 meses | ~16.666.667 |
| 36 meses | ~33.333.333 |
| 48 meses | 50.000.000 |

### Reglas del vesting

- Los tokens en vesting **NO se pueden transferir** hasta que se liberen.
- Los tokens en vesting **SÍ se pueden usar para staking** (delegar).
- Los tokens en vesting **SÍ cuentan para votar** en gobernanza.
- Las recompensas de staking son **completamente LÍQUIDAS** y no están sujetas a
  vesting.

---

## 4. Multisig de cuentas custodiadas

Las cuatro cuentas de la comunidad (treasury, airdrop, incentives, reserve) están
controladas por direcciones multisig.

| Cuenta | Tipo | Firmantes | Firmas necesarias |
|--------|------|-----------|-------------------|
| treasury | 3 de 3 | henry, george, charlie | Las 3 (unanimidad) |
| airdrop | 2 de 4 | henry, george, charlie, airdrop-aux | 2 de los 3 fundadores |
| incentives | 2 de 4 | henry, george, charlie, incentives-aux | 2 de los 3 fundadores |
| reserve | 2 de 4 | henry, george, charlie, reserve-aux | 2 de los 3 fundadores |

> **Nota sobre las claves auxiliares.** Las claves `airdrop-aux`,
> `incentives-aux` y `reserve-aux` existen **únicamente** para generar una
> dirección multisig distinta a la de treasury. **No participan en la firma real**:
> en la práctica, cada una de esas cuentas se opera con la firma de 2 de los 3
> fundadores.

---

## 5. Política de tesorería

Distribución orientativa (no técnica) de los 300.000.000 WORRELL custodiados en la
cuenta `treasury`.

| Partida | % de treasury | WORRELL |
|---------|--------------:|--------:|
| Desarrollo | 25 % | 75.000.000 |
| Liquidez | 20 % | 60.000.000 |
| Grants | 15 % | 45.000.000 |
| Partnerships | 15 % | 45.000.000 |
| Marketing | 15 % | 45.000.000 |
| Contingencias | 10 % | 30.000.000 |
| **TOTAL** | **100 %** | **300.000.000** |

---

## 6. Inflación y recompensas

El módulo de mint emite nuevos tokens en función del porcentaje del supply que esté
en staking, ajustando la inflación dinámicamente.

| Parámetro | Valor |
|-----------|-------|
| Inflación mínima | 7 % |
| Inflación máxima | 13 % |
| Inflación inicial | 13 % |
| Goal bonded ratio | 67 % |
| Mint denom | uworrell |

Cuando menos del 67 % del supply está en staking, la inflación sube hacia el 13 %;
cuando hay más, baja hacia el 7 %.

### Distribución de recompensas

| Destino | % de las recompensas |
|---------|---------------------:|
| Validadores y delegadores | 90 % |
| Community pool (community tax) | 10 % |

El community tax es del **10 %**: ese porcentaje de las recompensas de staking se
deriva al community pool y el **90 %** restante se reparte entre validadores y
delegadores.

---

## 7. Fees de transacción

| Parámetro | Valor |
|-----------|-------|
| Min gas price | 0,025 uworrell |

Cada nodo aplica un precio mínimo de gas de **0,025 uworrell**. Las transacciones
que ofrezcan un precio inferior serán rechazadas por el mempool.
