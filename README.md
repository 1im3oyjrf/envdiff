# envdiff

Compare `.env` files across environments and report missing or mismatched keys with optional secret masking.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
# Compare two .env files
envdiff .env.development .env.production

# Mask secret values in output
envdiff --mask .env.development .env.production

# Compare multiple environments against a base
envdiff --base .env.example .env.staging .env.production
```

### Example Output

```
MISSING in .env.production:
  - DATABASE_URL
  - REDIS_HOST

MISMATCHED keys:
  - APP_ENV  (development vs production)
  - LOG_LEVEL  (debug vs warning)
```

---

## Flags

| Flag | Description |
|------|-------------|
| `--mask` | Mask secret values in diff output |
| `--base` | Specify a base file to compare others against |
| `--json` | Output results in JSON format |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)