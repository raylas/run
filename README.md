# run

## Catalog

Place scripts in `catalog/`.

For positional argument discovery, use the following format for script headers:

```bash
#! /usr/bin/env nix-shell
#! nix-shell -i bash -p inetutils
# shellcheck shell=bash
#
## Traceroute to a host
# host: Host to traceroute to
# timeout: Timeout in seconds [30]
```

Where:

- `## Traceroute to a host` is the subcommand description
- `# timeout: Timeout in seconds [30]` is:
  - `timeout` is the flag name
  - `Timeout in seconds` is the flag description
  - `30` is the default value
