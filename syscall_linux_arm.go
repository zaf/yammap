//go:build linux && arm
// +build linux,arm

/*
	Copyright (C) 2022, Lefteris Zafiris <zaf@fastmail.com>
	This program is free software, distributed under the terms of
	the GNU GPL v3 License. See the LICENSE file
	at the top of the source tree.
*/

package yammap

const (
	SYS_MMAP      = 192
	SYS_MREMAP    = 163
	SYS_MUNMAP    = 91
	SYS_MSYNC     = 144
	SYS_FTRUNCATE = 93
	SYS_MADVISE   = 220

	maxSize = (1 << 31) - 1 // maximum allocation size, 2GiB for 32bit CPUs
)
