# Guía de contribución — Worrell

¡Gracias por tu interés en contribuir a Worrell! Esta guía describe el proceso para
proponer cambios al proyecto.

## Código de conducta

Esperamos un trato respetuoso y profesional entre todas las personas que participan.
No se tolera el acoso ni el comportamiento abusivo. Al participar, aceptas mantener
un entorno colaborativo y constructivo.

## Antes de empezar

- Lee el [README](README.md) y la documentación en [`docs/`](docs/) para entender la
  arquitectura, los parámetros de red y el modelo de gobernanza.
- Para cambios grandes o que afecten a parámetros de consenso/economía, **abre primero
  un issue** para discutir la propuesta. Muchos cambios de parámetros de la red en
  producción se realizan vía gobernanza on-chain, no por PR (ver [docs/GOVERNANCE.md](docs/GOVERNANCE.md)).

## Flujo de trabajo

1. **Fork** del repositorio `github.com/worrellchain/worrell`.
2. Crea una **rama** descriptiva a partir de la rama principal:
   ```bash
   git checkout -b feat/breve-descripcion
   ```
3. Realiza tus cambios con **commits** claros y atómicos. Usa mensajes en imperativo
   (p. ej. `fix: corrige cálculo de inflación`, `docs: amplía guía de validadores`).
   Se recomienda el estilo [Conventional Commits](https://www.conventionalcommits.org/).
4. Asegúrate de que el proyecto **compila** y las comprobaciones pasan:
   ```bash
   ignite chain build
   go test ./...
   ```
5. Abre un **Pull Request** contra la rama principal. Describe el qué y el porqué,
   enlaza los issues relacionados e incluye pasos de prueba.

## Requisitos para los PR

- El código compila (`ignite chain build`) y los tests pasan (`go test ./...`).
- Sigue el estilo del código existente (`gofmt`/`go vet` limpios).
- Cambios de comportamiento acompañados de tests cuando sea posible.
- Documentación actualizada si el cambio afecta a parámetros, comandos o flujos.
- **No** incluyas claves privadas, mnemónicos, archivos `.env` ni datos de nodo en el PR.

## Seguridad

Si encuentras una vulnerabilidad, **no abras un issue público**. Repórtala de forma
responsable a **security@worrellchain.io**. Consulta la sección de seguridad del
[README](README.md#seguridad) para más detalles.

## Licencia

Al contribuir, aceptas que tus aportaciones se publiquen bajo la licencia
[Apache 2.0](LICENSE) del proyecto.
