# Envcrypto

Proof of Concept

## Demo

```bash
datadir=$(mktemp -d)

# Generate an age keypair:
age-keygen -o "$datadir/age-key.txt"
age_public_key=$(age-keygen -y "$datadir/age-key.txt")

# Generate a new envcrypto environment:
envcrypto generate > "$datadir/envcrypto.env"

# Use sops to encrypt the private key with the age key:
sops --encrypt --age "$age_public_key" --encrypted-regex "PRIVATE_KEY" "$datadir/envcrypto.env" > "$datadir/envcrypto.sops.env"

# Encrypt a value with the envcrypto public key, and write it to a dotenv file:
echo "DOTENV_EXAMPLE=$(envcrypto -f "$datadir/envcrypto.sops.env" encrypt "hello dotenv")" > "$datadir/example.env"

# Decrypt the value by using the embedded support for sops:
SOPS_AGE_KEY_FILE="$datadir/age-key.txt" envcrypto --sops "$datadir/envcrypto.sops.env" -f "$datadir/example.env" get DOTENV_EXAMPLE

# Encrypt a value with the envcrypto public key, and export it:
export "ENV_EXAMPLE=$(envcrypto -f "$datadir/envcrypto.sops.env" encrypt "hello environment")"

# Decrypt the value by using the embedded support for sops:
SOPS_AGE_KEY_FILE="$datadir/age-key.txt" envcrypto --sops "$datadir/envcrypto.sops.env" get ENV_EXAMPLE

# Clean up:
rm -rf "$datadir"
```
