# hcelper
```
edmund@oncall:~/git/hcelper> ./hcelper login --help
Error: required flag(s) "username" not set
Usage:
  hcelper login [flags]

Flags:
  -e, --env string        The environment you're logging into (pre or prod)
  -h, --help              help for login
  -m, --method string     The login method (default "ldap")
  -n, --namespace string  The target namespace
  -u, --username string   The username for user login credentials

Global Flags:
      --config string   config file (default is .hce.hcl)

required flag(s) "username" not set
```

`hcelper` is designed to be used to aid with navigating the Hashicorp enterprise stack. The idea is to use a Vault-first approach that logs you into Vault and then uses secrets engines and configured roles to export the environmental variables required to access the rest of your stack.

For people working with multiple tenant spaces, this helper program provides a simple method to log in to ALL project-related namespaces by logging into just vault. This is accomplished by the fact that access to Hashicorp endpoints is generally accomplished by an `ADDR` field and a `TOKEN` field. This includes vault, though its default behaviour is to stick the token in a file (ew), which interferes with running multiple shells, as they all source the same file.

Early on I realized I'd forgotten basic unix shell philosophy, namely that child shells can't export information to a parent shell, so the initial idea that the helper would export that for you was nix'd. `eval`ing is out too because I need to prompt for authentication and I super do not want to tell people to just store that in a file. This might be achieveable by reading a var prefixing the command, but I also hate that because then the password goes into your shell history unless you've got `HISTIGNORE` set in your rc file.

So for now copy and paste the output!

## Requirements

- Vault Enterprise binary
- hcelper
- Nomad Enterprise binary (optional)
- Consul Enterprise binary (optional)

# Usage

## Basic
`hcelper login --username eadderley`

Currently this logs you into the root namespace and dumps `export` values of all your secret backends, like so:

```
export VAULT_ADDR=https://vault.url.systems:8200
export VAULT_TOKEN=NOPE
export CONSUL_HTTP_ADDR=https://consul.words.systems:8501
export CONSUL_HTTP_TOKEN=REDACTED
export NOMAD_ADDR=https://nomad.something.systems:4646
export NOMAD_TOKEN=NUFFIN
```

## Advanced 

`hcelper login -u eadderley -e pre -n lab-ce-shared`

Here you can actually target the approriate namespace.

In the future, you will be prompted to save a given set of selections (like env, namespace, secret backend roles), ordered by namespace which you will be able to access with:

`hcelper login -u eadderley -c $namespace`

## To Come

- config file support to persist selections and avoid repetitive menu navigation of secret backends.
- move config endpoints to config file fully to support full OSS release
- maybe finally provide a solution to automatically export this information