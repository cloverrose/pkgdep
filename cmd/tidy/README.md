# pkgdep_tidy

`pkgdep_tidy` is a tool to remove unused dependency rules from your configuration file.

## Overview

This tool provides:
- Detection of unused dependency rules
- Removal of unused rules from configuration
- Automatic cleanup of empty keys

## Installation

```shell
$ go install github.com/cloverrose/pkgdep/cmd/tidy@latest
```

### Or Build from source

```shell
$ make build/tidy
```

### Or Install via aqua

https://aquaproj.github.io/


## Usage

```bash
pkgdep_tidy -config=.pkgdep.yaml -inspector.file=used_rules.csv -output=.pkgdep.tidy.yaml
```

### Flags

- `-config` string Path to config file (default: `.pkgdep.yaml`)
  - The main configuration file containing dependency rules
- `-inspector.file` string Path to used rules record (default: `used_rules.csv`)
  - CSV file containing the actually used dependency rules
- `-output` string Path to output file (default: `.pkgdep.tidy.yaml`)
  - Where to write the tidied configuration
- `-log.level` string Log level (default: `info`)
  - Available levels: debug, info, warn, error
- `-log.format` string Log format (default: `json`)
  - Available formats: json, text

### Tips

For best results:

1. Format your input configuration file before running tidy:
   ```bash
   yq --prettyPrint --indent=2 .pkgdep.yaml
   ```

   This ensures:
   - Clean, consistent formatting
   - Meaningful diffs in the tidied output
   - Easier review of changes

2. Review the tidied configuration carefully before replacing your original file
