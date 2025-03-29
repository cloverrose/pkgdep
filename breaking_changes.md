# Breaking changes

This document lists breaking changes and deprecated features in pkgdep, along with migration guides to help users update their code.

## v0.4.0

### JSON configuration file

JSON configuration file is no longer supported. Please use YAML format instead.

### enableRegexp option

Regular expressions are now the only supported pattern matching syntax.

The `enableRegexp: false` configuration option has been removed. Attempting to use this option will result in an error.

**Migration Guide**

1. Remove the `enableRegexp` setting from your configuration file - regular expressions are now always enabled.

2. Update your pattern matching syntax:
    - Replace simple wildcards (`*`) with regular expressions:
        - Use `.*` for zero or more characters
        - Use `.?` for zero or one character
        - Use `.+` for one or more characters
    - To match an actual dot character, escape it with `\\.`
