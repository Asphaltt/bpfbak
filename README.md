<!--
 Copyright 2024 Leon Hwang.
 SPDX-License-Identifier: Apache-2.0
-->

# bpfbak: a tiny tool to backup bpf objects under bpffs

Currently, bpffs does not support `cp` pinned bpf objects. This tool is a workaround to backup pinned bpf objects under bpffs.

## Usage

`bpfbak` can be used as a library:

```shell
    opts := bpfbak.BackupOpts{
        Src:            rootFlags.src,
        UnpinSrc:       rootFlags.srcUnpin,
        Dst:            rootFlags.dst,
        AutoMountBpffs: rootFlags.autoMount,
        BpffsPath:      rootFlags.mountBpffs,
    }
    err := bpfbak.Backup(opts)
    if err != nil {
        log.Fatalf("Failed to backup bpf object: %v", err)
    }
```

Or as a command line tool:

```shell
$ ./bpfbak -h
bpfbak is a tool to backup eBPF objects

Usage:
  bpfbak [flags]

Flags:
      --auto-mount           automatically mount bpffs at the destination directory or --mount-bpffs
  -d, --dst string           destination filepath to backup the bpf object
  -h, --help                 help for bpfbak
      --mount-bpffs string   path to the directory where bpffs is mounted
  -s, --src string           source bpf object to be backed up
      --unpin-src            unpin the source bpf object after backing up
```

## License

Licensed under the Apache License, Version 2.0.
