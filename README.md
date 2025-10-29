# scc-product-version-verifier

A CLI tool to verify product versions in SUSE Customer Center.

`scc-product-version-verifier` is a command-line tool that interacts with the
SUSE Customer Center (SCC) API to verify if a specific product, version,
and architecture combination exists.

## Building

This project uses [GoReleaser](https://goreleaser.com/) to manage builds. To build the project, you can run:

```bash
goreleaser build --snapshot --clean
```

## Usage

The main command of this tool is `curl-verify`.

### `curl-verify`

This command mimics using curl to check if a product exists in SCC.

**Usage:**

```bash
scc-product-version-verifier curl-verify [product name] [product version] [product arch : optional] -R [registration code]
```

**Arguments:**

*   `[product name]`: The name of the product to verify.
*   `[product version]`: The version of the product to verify.
*   `[product arch]`: (Optional) The architecture of the product to verify. Defaults to `unknown`.

**Flags:**

*   `-R`, `--reg-code`: The SCC Registration Code used to authenticate for the API call.

**Example:**

```bash
scc-product-version-verifier curl-verify "sles" "15" "x86_64" -R "your-registration-code"
```

### `version`

Prints the version information of the CLI tool.

**Usage:**

```bash
scc-product-version-verifier --version
```

or

```bash
scc-product-version-verifier -V
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
