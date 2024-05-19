// Copyright 2024 Leon Hwang.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"

	"github.com/Asphaltt/bpfbak"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bpfbak",
	Short: "bpfbak is a tool to backup eBPF objects",
	Run: func(cmd *cobra.Command, args []string) {
		backupBpfObj()
	},
}

var rootFlags struct {
	src        string
	dst        string
	srcUnpin   bool
	autoMount  bool
	mountBpffs string
}

func init() {
	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&rootFlags.src, "src", "s", "", "source bpf object to be backed up")
	flags.StringVarP(&rootFlags.dst, "dst", "d", "", "destination filepath to backup the bpf object")
	flags.BoolVar(&rootFlags.srcUnpin, "unpin-src", false, "unpin the source bpf object after backing up")
	flags.BoolVar(&rootFlags.autoMount, "auto-mount", false, "automatically mount bpffs at the destination directory or --mount-bpffs")
	flags.StringVar(&rootFlags.mountBpffs, "mount-bpffs", "", "path to the directory where bpffs is mounted")
}

func main() {
	_ = rootCmd.Execute()
}

func backupBpfObj() {
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
}
