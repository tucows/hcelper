# hcelper

`hcelper` is designed to be used to aid with navigating the Hashicorp enterprise stack. The idea is to use a Vault-first approach that logs you into vault and then sets uses secrets engines and configured roles to export the environmental variables required to access the rest of your stack.

## Requirements

- Vault Enterprise binary
- hcelper
- Nomad Enterprise binary (optional)
- Consul Enterprise binary (optional)

# Usage

`hcelper login --env pre --username eadderley`