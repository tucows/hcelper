# hcelper

`hcelper` is designed to be used to aid with navigating the Hashicorp enterprise stack. The idea is to use a Vault-first approach that logs you into vault and then sets uses secrets engines and configured roles to export the environmental variables required to access the rest of your stack.

```
edmund@oncall:~/git/hcelper> ./hcelper login --help
Error: required flag(s) "username" not set
Usage:
  hcelper login [flags]

Flags:
  -e, --env string        The environment you're logging into (pre or prod)
  -h, --help              help for login
  -m, --method string     The login method (default "ldap")
  -u, --username string   The username for user login credentials

Global Flags:
      --config string   config file (default is .hce.hcl)

required flag(s) "username" not set
```

## Requirements

- Vault Enterprise binary
- hcelper
- Nomad Enterprise binary (optional)
- Consul Enterprise binary (optional)

# Usage

## Basic (working)
`hcelper login --username eadderley`

Currently this logs you into the root namespace and dumps `export` values of interest to the HCE team, like so:

```
export VAULT_ADDR=https://vault.url.systems:8200
export VAULT_TOKEN=NOPE
export CONSUL_HTTP_ADDR=https://consul.words.systems:8501
export CONSUL_HTTP_TOKEN=REDACTED
export NOMAD_ADDR=https://nomad.something.systems:4646
export NOMAD_TOKEN=NUFFIN
```

This will shortly receive two updates:
1. Set non-standard vault endpoint in config file (currently uses openstack's)
2. Ability to specify namespace (for everyone else)
3. Ability to select which roles you want to generate creds off of (for tenants with customized credential endpoints)

## Advanced (to come)

1. Export your answers to the prompts to a config file, ordered by namespace select from your pre-generated answers with a flag, simply log into vault
2. Support for other backends, like ssh

## FINAL FORM

????