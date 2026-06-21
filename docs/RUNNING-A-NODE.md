# Guía para validadores: ejecutar un nodo Worrell

Esta guía describe cómo poner en marcha un nodo completo de **Worrell** y cómo convertirlo en validador en la testnet. Worrell es una blockchain proof-of-stake construida sobre Cosmos SDK.

| Concepto | Valor |
|----------|-------|
| Binario | `worrelld` |
| Token | WORRELL (denominación base: `uworrell`, 6 decimales) |
| 1 WORRELL | 1.000.000 uworrell |
| Data directory | `~/.worrell` |
| Chain ID testnet | `worrell-testnet-1` |
| Chain ID mainnet | `worrell-1` |
| Min gas price | `0.025uworrell` |

---

## 1. Requisitos de hardware

### Testnet

| Recurso | Recomendación |
|---------|---------------|
| Sistema operativo | Ubuntu 22.04 LTS |
| CPU | 2 vCPU |
| RAM | 4 GB |
| Disco | 100 GB SSD |
| Red | Conexión estable, IP pública para los puertos P2P |

### Mainnet

Para un validador de mainnet se recomienda más holgura, ya que el nodo debe firmar bloques de forma ininterrumpida y mantener historial:

| Recurso | Recomendación |
|---------|---------------|
| Sistema operativo | Ubuntu 22.04 LTS |
| CPU | 4 vCPU o más |
| RAM | 16 GB o más |
| Disco | 500 GB+ NVMe SSD (con margen de crecimiento) |
| Red | Conexión estable de baja latencia, IP pública dedicada |

> El bloque se produce cada ~5-6 segundos. Un disco lento o una red inestable provocan bloques perdidos y, con el tiempo, *jailing* por *downtime*.

---

## 2. Instalación

### Dependencias

- **Go 1.22+**
- **git**, **make**, **build-essential**

Instalación de dependencias en Ubuntu 22.04:

```bash
sudo apt update && sudo apt install -y git curl build-essential
```

Instalación de Go (ejemplo con Go 1.22):

```bash
curl -LO https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.profile
source ~/.profile
go version
```

### Compilar desde código fuente

```bash
git clone https://github.com/worrellchain/worrell.git
cd worrell
make install
```

`make install` compila el binario `worrelld` y lo deja en `$HOME/go/bin`. Alternativamente, puedes compilar con Ignite CLI:

```bash
ignite chain build
```

Verifica la instalación:

```bash
worrelld version
worrelld --help
```

---

## 3. Inicialización del nodo

Inicializa la configuración local del nodo eligiendo un *moniker* (nombre público de tu nodo):

```bash
worrelld init <moniker> --chain-id worrell-testnet-1
```

Esto crea el directorio de datos en `~/.worrell` con la siguiente estructura relevante:

- `~/.worrell/config/genesis.json` — genesis (se reemplaza en el paso siguiente)
- `~/.worrell/config/config.toml` — configuración de CometBFT (peers, P2P, RPC)
- `~/.worrell/config/app.toml` — configuración de la aplicación (gas, API, gRPC)
- `~/.worrell/config/priv_validator_key.json` — clave de firma del validador (¡protégela!)
- `~/.worrell/config/node_key.json` — identidad de red del nodo

---

## 4. Unirse a la testnet

### 4.1 Obtener el genesis

Sustituye el genesis generado localmente por el genesis oficial de la red:

```bash
curl -s https://raw.githubusercontent.com/worrellchain/worrell/main/networks/testnet/genesis.json \
  -o ~/.worrell/config/genesis.json
worrelld genesis validate-genesis
```

### 4.2 Configurar peers y seeds

Edita `~/.worrell/config/config.toml` y configura `seeds` y `persistent_peers` en la sección `[p2p]` con los nodos publicados para la testnet (formato `<node_id>@<host>:26656`):

```toml
# ~/.worrell/config/config.toml  -> [p2p]
seeds = "<seed_node_id>@<seed_host>:26656"
persistent_peers = "<peer_node_id_1>@<peer_host_1>:26656,<peer_node_id_2>@<peer_host_2>:26656"
```

### 4.3 Fijar el precio mínimo de gas

Edita `~/.worrell/config/app.toml` y establece el precio mínimo de gas. Esto es obligatorio para que el nodo acepte y procese transacciones:

```toml
# ~/.worrell/config/app.toml
minimum-gas-prices = "0.025uworrell"
```

### 4.4 Arrancar el nodo

```bash
worrelld start
```

Se recomienda ejecutarlo como servicio `systemd` para que arranque automáticamente y se reinicie ante fallos. Ejemplo de unidad en `/etc/systemd/system/worrelld.service`:

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

### 4.5 Esperar a la sincronización

Antes de crear un validador, el nodo debe estar completamente sincronizado:

```bash
worrelld status 2>&1 | jq '.sync_info'
```

Cuando `catching_up` sea `false`, el nodo está al día. La red tiene **State Sync habilitado**, lo que permite sincronizar rápidamente sin descargar todo el historial.

---

## 5. Crear un validador

> Necesitas: nodo sincronizado (`catching_up: false`) y una cuenta con saldo suficiente en WORRELL. En testnet puedes pedir fondos al faucet.

### 5.1 Crear o importar una clave

```bash
# Crear una clave nueva (guarda la frase mnemónica en lugar seguro)
worrelld keys add <nombre-clave>

# O importar una existente
worrelld keys add <nombre-clave> --recover
```

Consulta la dirección y el saldo:

```bash
worrelld keys show <nombre-clave> -a
worrelld query bank balances $(worrelld keys show <nombre-clave> -a)
```

### 5.2 Obtener la clave pública de consenso

```bash
worrelld tendermint show-validator
```

Devuelve un valor con formato `{"@type":"/cosmos.crypto.ed25519.PubKey","key":"..."}` que usarás en `--pubkey`.

### 5.3 Comando `create-validator`

El comando moderno de Cosmos SDK toma un archivo JSON con la definición del validador. Crea `validator.json`:

```json
{
  "pubkey": {"@type":"/cosmos.crypto.ed25519.PubKey","key":"<TU_PUBKEY>"},
  "amount": "20000000000000uworrell",
  "moniker": "<moniker>",
  "identity": "",
  "website": "",
  "security": "",
  "details": "Validador de la testnet de Worrell",
  "commission-rate": "0.05",
  "commission-max-rate": "0.25",
  "commission-max-change-rate": "0.01",
  "min-self-delegation": "1000000"
}
```

Y envíalo:

```bash
worrelld tx staking create-validator validator.json \
  --from <nombre-clave> \
  --chain-id worrell-testnet-1 \
  --gas auto \
  --gas-adjustment 1.5 \
  --gas-prices 0.025uworrell \
  --yes
```

Si tu versión del binario usa la forma con flags, el equivalente es:

```bash
worrelld tx staking create-validator \
  --pubkey="$(worrelld tendermint show-validator)" \
  --amount="20000000000000uworrell" \
  --moniker="<moniker>" \
  --commission-rate="0.05" \
  --commission-max-rate="0.25" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1000000" \
  --from=<nombre-clave> \
  --chain-id=worrell-testnet-1 \
  --gas=auto \
  --gas-adjustment=1.5 \
  --gas-prices=0.025uworrell \
  --yes
```

### 5.4 Explicación de los valores

| Flag / campo | Valor del ejemplo | Significado |
|--------------|-------------------|-------------|
| `--amount` | `20000000000000uworrell` | Autodelegación inicial (en este ejemplo, 20.000.000 WORRELL). Ajusta a tu caso. |
| `--commission-rate` | `0.05` | Comisión actual: 5%. **No puede ser inferior al mínimo global de red (5%).** |
| `--commission-max-rate` | `0.25` | Comisión máxima que **este validador** se permite a sí mismo: 25%. La fijas tú. |
| `--commission-max-change-rate` | `0.01` | Cambio máximo de comisión por día: 1%. Lo fijas tú. |
| `--min-self-delegation` | `1000000` | Autodelegación mínima que mantendrás. **`1000000` uworrell = 1 WORRELL.** |

> ⚠️ **Crítico — `--min-self-delegation`:** el valor se expresa en `uworrell`, no en WORRELL. Para garantizar el mínimo de red de 1 WORRELL debes usar **`--min-self-delegation="1000000"`**. **NUNCA** uses `"1"`: eso significaría 1 uworrell (0,000001 WORRELL) y dejaría tu validador por debajo del mínimo efectivo previsto.

### 5.5 Comisiones: mínimo global vs. máximo propio

- El **`min_commission_rate` del 5% (0.05) es GLOBAL**: es el mínimo que impone la red. Ningún validador puede fijar una `commission-rate` por debajo de 0.05.
- En cambio, el **máximo de comisión** (`commission-max-rate`) y el **cambio máximo diario** (`commission-max-change-rate`) los decide **cada validador por su cuenta**. Una vez fijados al crear el validador, el `commission-max-rate` **no se puede aumentar después**, así que elígelo con cuidado.

### 5.6 Verificar que el validador está activo

```bash
worrelld query staking validator $(worrelld keys show <nombre-clave> --bech val -a)
```

Comprueba que el estado sea `BOND_STATUS_BONDED` y que `jailed` sea `false`.

---

## 6. Monitorización

### 6.1 Estado del nodo

```bash
worrelld status
worrelld status 2>&1 | jq '.sync_info.latest_block_height, .sync_info.catching_up'
```

### 6.2 Uptime y bloques firmados

Consulta la información de firma de tu validador (necesitas su *consensus address*):

```bash
# Dirección de consenso del validador
worrelld query slashing signing-info $(worrelld tendermint show-address)
```

El resultado incluye `missed_blocks_counter` (bloques no firmados dentro de la ventana) y `jailed_until`. Vigila que `missed_blocks_counter` se mantenga bajo.

### 6.3 Parámetros de slashing relevantes (evitar el jailing)

| Parámetro | Valor | Qué significa |
|-----------|-------|---------------|
| Ventana de bloques firmados | 10.000 bloques | Periodo en el que se mide la disponibilidad |
| Mínimo firmado por ventana | 5% (0.05) | Debes firmar al menos el 5% de los bloques de la ventana |
| Penalización por downtime | 0,01% (0.0001) | Slash al ser encarcelado por inactividad |
| Penalización por doble firma | 5% (0.05) | Slash por firmar dos bloques al mismo height |
| Duración del jail por downtime | 300s (5 min) en testnet / 43.200s (12 h) en mainnet | Tiempo mínimo encarcelado |

Para **evitar el jailing por downtime**:

- Ejecuta el nodo como servicio `systemd` con reinicio automático.
- Mantén disco y CPU con holgura; un nodo que no sigue el ritmo pierde bloques.
- **Nunca** ejecutes dos instancias con la misma `priv_validator_key.json` a la vez: provocaría una **doble firma** y un slash del 5% (mucho más grave que un downtime).
- Monitoriza con Prometheus (ver abajo) y configura alertas sobre bloques perdidos.

### 6.4 Recuperar un validador encarcelado (unjail)

Si tu validador fue encarcelado por downtime, una vez transcurrido el periodo de jail puedes reactivarlo:

```bash
worrelld tx slashing unjail \
  --from <nombre-clave> \
  --chain-id worrell-testnet-1 \
  --gas auto \
  --gas-adjustment 1.5 \
  --gas-prices 0.025uworrell \
  --yes
```

### 6.5 Prometheus

La telemetría Prometheus está habilitada. Para exponer métricas, en `~/.worrell/config/config.toml`:

```toml
# ~/.worrell/config/config.toml  -> [instrumentation]
prometheus = true
prometheus_listen_addr = ":26660"
```

Las métricas quedan disponibles en `http://<host>:26660/metrics`.

---

## 7. Puertos

Asegúrate de que el cortafuegos permite el tráfico necesario. El puerto P2P debe ser accesible desde el exterior; el resto, preferiblemente, restringido a localhost o a IPs de confianza.

| Servicio | Puerto | Exposición recomendada |
|----------|--------|------------------------|
| CometBFT P2P | 26656 | Pública (necesaria para conectarse a la red) |
| CometBFT RPC | 26657 | Localhost / IPs de confianza |
| API REST | 1317 | Localhost por defecto |
| gRPC | 9090 | Localhost por defecto |
| Prometheus | 26660 | Localhost / red interna de monitorización |

Por defecto, la API REST y gRPC escuchan solo en localhost. Si los expones públicamente, protégelos con un proxy inverso y reglas de cortafuegos.

---

## 8. Lista de comprobación rápida

- [ ] Go 1.22+ instalado y `worrelld version` funciona
- [ ] Nodo inicializado con `--chain-id worrell-testnet-1`
- [ ] `genesis.json` oficial colocado y validado
- [ ] `seeds` / `persistent_peers` configurados
- [ ] `minimum-gas-prices = "0.025uworrell"` en `app.toml`
- [ ] Nodo sincronizado (`catching_up: false`)
- [ ] Cuenta con saldo suficiente
- [ ] Validador creado con `--min-self-delegation="1000000"` (NO `"1"`)
- [ ] `commission-rate` ≥ 0.05 (mínimo global de red)
- [ ] Nodo corriendo bajo `systemd` con reinicio automático
- [ ] Monitorización de bloques firmados activa
