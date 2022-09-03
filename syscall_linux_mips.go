//go:build linux && mips
// +build linux,mips

/*
	Copyright (C) 2022, Lefteris Zafiris <zaf@fastmail.com>
	This program is free software, distributed under the terms of
	the GNU GPL v3 License. See the LICENSE file
	at the top of the source tree.
*/

package yammap

const (
	SYS_MMAP      = 4090
	SYS_MREMAP    = 4167
	SYS_MUNMAP    = 4091
	SYS_MSYNC     = 4144
	SYS_FTRUNCATE = 4212
	SYS_MADVISE   = 4218

	maxSize = 1 << 31 // maximum allocation size, 2GiB for 32bit CPUs
)
