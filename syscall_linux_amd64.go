//go:build linux && amd64
// +build linux,amd64

/*
	Copyright (C) 2022, Lefteris Zafiris <zaf@fastmail.com>
	This program is free software, distributed under the terms of
	the GNU GPL v3 License. See the LICENSE file
	at the top of the source tree.
*/

package yammap

const (
	SYS_MMAP      = 9
	SYS_MREMAP    = 25
	SYS_MUNMAP    = 11
	SYS_MSYNC     = 26
	SYS_FTRUNCATE = 77
	SYS_MADVISE   = 28

	maxSize = (1 << 47) - 1 // maximum allocation size, 128TiB for x86_64
)
