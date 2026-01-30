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

*   `-R`, `--regcode`: The SCC Registration Code used to authenticate for the API call. Can also be set via the `SCC_REGCODE` environment variable.

**Example:**

```bash
scc-product-version-verifier curl-verify "SLES" "15" "x86_64" -R "your-registration-code"
```

You can also use the `SCC_REGCODE` environment variable (better for CI):
```bash
export SCC_REGCODE="your-registration-code"
scc-product-version-verifier curl-verify "SLES" "15" "x86_64"
```

```bash
export SCC_REGCODE="your-registration-code"
scc-product-version-verifier curl-verify rancher 2.12.3
```

> [!WARNING]
> The SCC api is case-sensitive for product lookup meaning `SLES` != `sles`.
> For SLES look up it must be upper case, for `rancher` lookup it must be lower case.

## GitHub Actions

This repository provides reusable GitHub Actions to download and use the verifier in your CI/CD workflows.

### Download Action

Downloads and installs the latest version of `scc-product-version-verifier`.

**Location:** `rancherlabs/scc-product-version-verifier/actions/download`

**Requirements:**
- Works on Linux runners
- No sudo required (uses GitHub Actions provided gh cli)

**Outputs:**
- `version`: The installed version of the verifier
- `bin-path`: Installation path of the verifier
- `asset-name`: Name of the downloaded tool

**Example:**

```yaml
- name: Setup SCC Product Version Verifier
  uses: rancherlabs/scc-product-version-verifier/actions/download@main
```

### Verify Action

Verifies a product version against SCC staging and/or production environments.

**Location:** `rancherlabs/scc-product-version-verifier/actions/verify`

**Requirements:**
- `scc-product-version-verifier` must be installed (use the download action first)
- Valid SCC registration code(s)

**Inputs:**
- `version` (required): Version to verify (will be sanitized to remove `v` prefix and prerelease suffixes)
- `staging-code` (optional): SCC staging registration code
- `production-code` (optional): SCC production registration code
- `product-name` (required): Product name to verify (case-sensitive)
- `skip-staging` (optional, default: `false`): Skip staging verification
- `skip-production` (optional, default: `false`): Skip production verification
- `fail-on-error` (optional, default: `false`): Fail the workflow if verification fails

**Outputs:**
- `staging-result`: Staging verification result (`passed`/`failed`/`skipped`)
- `production-result`: Production verification result (`passed`/`failed`/`skipped`)

**Example:**

```yaml
- name: Setup Verifier
  uses: rancherlabs/scc-product-version-verifier/actions/download@main

- run: echo "${{ github.workspace }}/bin" >> $GITHUB_PATH

- name: Verify Product Version
  uses: rancherlabs/scc-product-version-verifier/actions/verify@main
  with:
    version: v2.12.3
    staging-code: ${{ secrets.SCC_STAGING_CODE }}
    production-code: ${{ secrets.SCC_PRODUCTION_CODE }}
    product-name: rancher
    fail-on-error: false
```

**Notes:**
- By default, verification failures do NOT fail the workflow (`fail-on-error: false`). Set to `true` to enforce strict verification.
- Version strings are automatically sanitized (e.g., `v2.12.3-rc1` becomes `2.12.3`)
- Product names are case-sensitive (e.g., `SLES` vs `sles`, `rancher` vs `Rancher`)
- Results are written to the GitHub Actions step summary for easy viewing