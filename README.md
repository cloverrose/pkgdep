# pkgdep

`pkgdep` checks if package dependency follows rule.

## Install

```shell
$ go install github.com/cloverrose/pkgdep/cmd/pkgdep@latest
```

## Or Build from source

```shell
$ go build -o bin/ ./cmd/...
```

## Usage

### 1. Create .pkgdep.json in your repository

See [.pkgdep.json](./.pkgdep.json) as example.

We can use `*` as wild card.

### 2. Run

```shell
$ go vet -vettool=`which pkgdep` -pkgdep.config=$(PWD)/.pkgdep.json ./...
```
