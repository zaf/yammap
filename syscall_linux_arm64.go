//go:build linux && arm64
// +build linux,arm64

/*
	Copyright (C) 2022, Lefteris Zafiris <zaf@fastmail.com>
	This program is free software, distributed under the terms of
	the GNU GPL v3 License. See the LICENSE file
	at the top of the source tree.
*/

package yammap

const (
	SYS_MMAP   = 222
	SYS_MREMAP = 215
	SYS_MUNMAP = 216
	SYS_MSYNC  = 227

	maxSize = 0xFFFFFFFFFFFF // maximum allocation size, 2^48 bytes for arm64
)
