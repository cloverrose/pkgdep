# pkgdep

`pkgdep` checks if package dependency follows rule.

## Config

`pkgdep` requires you to specify a configuration file using the `config` option.

You can configure via commandline option or golangci setting.

## Install

```shell
$ go install github.com/cloverrose/pkgdep/cmd/pkgdep@latest
```

### Or Build from source

```shell
$ make build
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
    import: 'github.com/cloverrose/mockguard'
    version: v0.4.0
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
```
