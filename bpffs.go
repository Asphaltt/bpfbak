// Copyright 2024 Leon Hwang.
// SPDX-License-Identifier: Apache-2.0

package bpfbak

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iovisor/gobpf/pkg/bpffs"
	"golang.org/x/sys/unix"
)

func isBpffsDir(filepath string) (bool, error) {
	var stat unix.Statfs_t

	err := unix.Statfs(filepath, &stat)
	if err != nil {
		return false, fmt.Errorf("failed to statfs %s: %v", filepath, err)
	}

	return stat.Type == unix.BPF_FS_MAGIC, nil
}

func isInBpffs(dirpath string) (bool, error) {
	if ok, _ := isBpffsDir(dirpath); ok {
		return true, nil
	}

	parent := filepath.Dir(dirpath)
	if parent == dirpath {
		return false, fmt.Errorf("failed to find bpffs in %s", dirpath)
	}

	return isInBpffs(parent)
}

func mountBpffs(dirpath string) error {
	err := os.MkdirAll(dirpath, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create dir %s: %v", dirpath, err)
	}

	err = bpffs.MountAt(dirpath)
	if err != nil {
		return fmt.Errorf("failed to mount bpffs at %s: %v", dirpath, err)
	}

	return nil
}

func autoMountBpffs(bpffsPath string, autoMount bool) error {
	if !autoMount {
		return fmt.Errorf("bpffs is not mounted at %s", bpffsPath)
	}

	if err := mountBpffs(bpffsPath); err != nil {
		return err
	}

	return nil
}

func prepareBpffs(bpffsPath string, autoMount bool) error {
	if inBpffs, _ := isInBpffs(bpffsPath); inBpffs {
		return nil
	}

	return autoMountBpffs(bpffsPath, autoMount)
}
