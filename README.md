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

### 1. Create .pkgdep.json in your repository

See [.pkgdep.json](./.pkgdep.json) as example.

- When `enableRegexp = false`: We can use `*` as wild card.
- When `enableRegexp = true`: We can use regexp.

### 2. Run

#### A. Use as go vet tool

```shell
$ go vet -vettool=`which pkgdep` -pkgdep.config=$(PWD)/.pkgdep.json ./...
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
    version: v0.3.5
```

`.golangci.yml`

```yaml
linters-settings:
  custom:
    pkgdep:
      type: "module"
      description: pkgdep validates if package dependency follows rule.
      settings:
        config: "/path/to/.pkgdep.yaml"
```
