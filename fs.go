// Copyright 2024 Leon Hwang.
// SPDX-License-Identifier: Apache-2.0

package bpfbak

import (
	"os"

	"golang.org/x/sys/unix"
)

func isexists(filepath string) bool {
	var stat unix.Stat_t
	err := unix.Stat(filepath, &stat)
	return err == nil || !os.IsNotExist(err)
}
