# Contribution Guide — Worrell

Thank you for your interest in contributing to Worrell! This guide describes the process for
proposing changes to the project.

## Code of Conduct

We expect respectful and professional treatment among everyone who participates.
Harassment and abusive behavior are not tolerated. By participating, you agree to maintain
a collaborative and constructive environment.

## Before You Start

- Read the [README](README.md) and the documentation in [`docs/`](docs/) to understand the
  architecture, network parameters, and governance model.
- For large changes or those that affect consensus/economic parameters, **open an issue
  first** to discuss the proposal. Many parameter changes to the network in
  production are made via on-chain governance, not via PR (see [docs/GOVERNANCE.md](docs/GOVERNANCE.md)).

## Workflow

1. **Fork** the `github.com/worrellchain/worrell` repository.
2. Create a descriptive **branch** from the main branch:
   ```bash
   git checkout -b feat/short-description
   ```
3. Make your changes with clear, atomic **commits**. Use imperative messages
   (e.g. `fix: correct inflation calculation`, `docs: expand validator guide`).
   The [Conventional Commits](https://www.conventionalcommits.org/) style is recommended.
4. Make sure the project **compiles** and the checks pass:
   ```bash
   ignite chain build
   go test ./...
   ```
5. Open a **Pull Request** against the main branch. Describe the what and the why,
   link the related issues, and include testing steps.

## PR Requirements

- The code compiles (`ignite chain build`) and the tests pass (`go test ./...`).
- It follows the style of the existing code (clean `gofmt`/`go vet`).
- Behavioral changes accompanied by tests whenever possible.
- Documentation updated if the change affects parameters, commands, or flows.
- Do **not** include private keys, mnemonics, `.env` files, or node data in the PR.

## Security

If you find a vulnerability, **do not open a public issue**. Report it
responsibly to **security@worrellchain.com**. See the security section of the
[README](README.md#security) for more details.

## License

By contributing, you agree that your contributions are released under the project's
[Apache 2.0](LICENSE) license.
