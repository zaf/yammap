//go:build linux && mips64le
// +build linux,mips64le

/*
	Copyright (C) 2022, Lefteris Zafiris <zaf@fastmail.com>
	This program is free software, distributed under the terms of
	the GNU GPL v3 License. See the LICENSE file
	at the top of the source tree.
*/

package yammap

const (
	SYS_MMAP      = 5009
	SYS_MREMAP    = 5024
	SYS_MUNMAP    = 5011
	SYS_MSYNC     = 5025
	SYS_FTRUNCATE = 5075
	SYS_MADVISE   = 5027

	maxSize = (1 << 47) - 1 // maximum allocation size, 128TiB for 64bit CPUs
)
