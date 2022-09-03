//go:build linux && riscv
// +build linux,riscv

/*
	Copyright (C) 2022, Lefteris Zafiris <zaf@fastmail.com>
	This program is free software, distributed under the terms of
	the GNU GPL v3 License. See the LICENSE file
	at the top of the source tree.
*/

package yammap

const (
	SYS_MMAP      = 222
	SYS_MREMAP    = 216
	SYS_MUNMAP    = 215
	SYS_MSYNC     = 227
	SYS_FTRUNCATE = 46
	SYS_MADVISE   = 233

	maxSize = 1 << 31 // maximum allocation size, 2GiB for 32bit CPUs
)
