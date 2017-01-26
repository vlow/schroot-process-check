# schroot-process-check

A simple command line tool to check if there are any processes running in a [schroot](https://packages.debian.org/sid/schroot) session.
Created to check if a schroot session can be safely ended.

To use this as a normal user, setguid and setgid can be used.
Assuming the binary's location is `/usr/bin/schroot-process-check`.

```bash
# chown root:wheel /usr/bin/schroot-process-check
# chmod 6711 /usr/bin/schroot-process-check
```

## Usage

This is the output of the tool when called without any argument:

```
Usage: schroot-process-check [OPTION]... SCHROOT-SESSION-NAME
Options:
  -q	Quiet mode, avoid all output.
  -v	Verbose mode, prints IDs of processes running in the given schroot session.
```

Additionally, the result can be gathered from the tools returncode:

| Active Processes Found | Returncode |
| ---------------------- | ---------- |
| no                     | 0          |
| yes                    | 3          |
