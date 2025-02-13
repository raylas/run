# run

Run autosemantic scripts in a Kubernetes context, or locally.

```bash
run [command]
```

## Catalog

Place scripts in `catalog/`.

For positional argument discovery, use the following format for script headers:

```bash
#! /usr/bin/env nix-shell
#! nix-shell -i bash -p inetutils
#
## Traceroute to a host
# host: Host to traceroute to
# timeout: Timeout in seconds [30]
```

Where:

- `#! /usr/bin/env nix-shell` let's Nix know what's coming
- `#! nix-shell -i bash -p inetutils` Tells nix-shell to:
  - Use `bash` as the shell
  - Make `inetutils` available in said shell
- `## Traceroute to a host` is the subcommand description
- `# timeout: Timeout in seconds [30]` is:
  - `timeout` is the flag name
  - `Timeout in seconds` is the flag description
  - `30` is the default value
