# Security Policy

## Reporting a vulnerability

If you discover a security vulnerability in Worrell — the node software, the network
configuration, or any related infrastructure — please report it **privately**:

**security@worrellchain.com**

Please do **not** open a public GitHub issue, discuss it in public channels, or
exploit it on any network.

### What to include

- A description of the vulnerability and its potential impact.
- Steps to reproduce (a proof of concept helps, but is not required).
- The affected component and version/commit, if known.
- How you would like to be credited, if you wish to be.

### What to expect

- **Acknowledgement within 72 hours** of your report.
- We will investigate, keep you informed of progress, and work on a fix.
- Please allow us a reasonable window to release a fix before any public
  disclosure. We will coordinate the disclosure timeline with you.

## Scope

- `worrelld` node software (this repository).
- Network configuration and genesis files ([worrellchain/networks](https://github.com/worrellchain/networks)).
- Public network infrastructure (testnet nodes, faucet).

Testnet tokens have no value; issues that only allow acquiring testnet funds are
appreciated but low severity. Consensus, key-management, fund-safety and
remote-code-execution issues are the highest priority.

## Supported versions

Only the latest release is supported. Please verify the issue is present in the
most recent version before reporting.
