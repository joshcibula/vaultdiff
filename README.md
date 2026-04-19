# vaultdiff

> CLI tool to diff secrets between two HashiCorp Vault paths or namespaces

---

## Installation

```bash
go install github.com/yourusername/vaultdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultdiff.git
cd vaultdiff
go build -o vaultdiff .
```

---

## Usage

```bash
# Diff secrets between two Vault paths
vaultdiff secret/prod/app secret/staging/app

# Diff across namespaces
vaultdiff --namespace-a team-a --namespace-b team-b secret/config

# Output in JSON format
vaultdiff --format json secret/prod/app secret/staging/app
```

**Example output:**

```
~ db/password   [changed]
+ feature_flag  [only in secret/prod/app]
- debug_mode    [only in secret/staging/app]
```

---

## Configuration

`vaultdiff` respects standard Vault environment variables:

| Variable | Description |
|---|---|
| `VAULT_ADDR` | Vault server address |
| `VAULT_TOKEN` | Authentication token |
| `VAULT_NAMESPACE` | Default namespace |

---

## Requirements

- Go 1.21+
- HashiCorp Vault 1.x
- A valid Vault token with read access to compared paths

---

## License

MIT © 2024 [Your Name](https://github.com/yourusername)