# pkgdep

`pkgdep` checks if package dependency follows rule.

## Options

- `config` (required)
  - Path to the configuration file (e.g. `.pkgdep.yaml`)
  - Supported file extensions: `.yaml`, `.yml`
- `log.level` (optional)
  - Controls logging verbosity
  - Valid values: `DEBUG`, `INFO`, `WARN`, `ERROR`
  - Default: `INFO`
- `log.file` (optional)
  - Path to write log output
  - If unspecified, logs to stdout
  - Note: When using with golangci-lint, a log file must be specified since stdout is captured
- `log.format` (optional)
  - Controls log output format
  - Valid values: `json`, `text`
  - Default: `json`
- `inspector.file` (optional, experimental)
  - Path to write inspector output
  - If unspecified, inspector output is not saved
  - Inspector output contains information about which dependency rules were used
  - **Warning**: The inspector feature is experimental and may change or be removed in future versions without notice.

You can configure via commandline option or golangci setting.

## Install

```shell
$ go install github.com/cloverrose/pkgdep/cmd/pkgdep@latest
```

### Or Build from source

```shell
$ make build/pkgdep
```

### Or Install via aqua

https://aquaproj.github.io/

## Usage

### 1. Create .pkgdep.yaml in your repository

See [.pkgdep.yaml](./.pkgdep.yaml) as example.

We can use regexp.

**Detailed Information**

`dependencies` is unmarshalled into an ordered map. Package dependencies are validated in order, starting from the first entry.

### 2. Run

#### A. Use as go vet tool

config file path should be absolute.

```shell
$ go vet -vettool=`which pkgdep` -pkgdep.config=$(PWD)/.pkgdep.yaml ./...
```

#### B. Use as golangci-lint custom plugin

https://golangci-lint.run/plugins/module-plugins/

Here are reference settings

`.custom-gcl.yml`

```yaml
version: v1.64.8
name: custom-golangci-lint
destination: bin
plugins:
  - module: 'github.com/cloverrose/pkgdep'
    import: 'github.com/cloverrose/pkgdep'
    version: v0.5.0
```

`.golangci.yml`

config file path can be relative.

```yaml
linters-settings:
  custom:
    pkgdep:
      type: "module"
      description: pkgdep validates if package dependency follows rule.
      settings:
        config: "./.pkgdep.yaml"
        log:
          level: "INFO"
          file: "./log.txt"
          format: "json"
        inspector:
          file: "./used_rules.csv"
```

# Other documents

## pkgdep_tidy

`pkgdep_tidy` is a command-line tool that helps clean up and optimize your `.pkgdep.yaml` configuration by:

- Removing unused dependency rules
- Cleaning up empty configuration keys
- Maintaining clean, consistent formatting

For detailed usage instructions and examples, see the [pkgdep_tidy documentation](cmd/tidy/README.md).

## pkgdep-tidy

`pkgdep-tidy` is a command-line tool that helps clean up and optimize your `.pkgdep.yaml` configuration by:

- Removing unused dependency rules
- Cleaning up empty configuration keys
- Maintaining clean, consistent formatting

For detailed usage instructions and examples, see the [pkgdep-tidy documentation](cmd/pkgdep-tidy/README.md).
